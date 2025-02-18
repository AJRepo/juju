// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	vault "github.com/mittwald/vaultgo"
	"gopkg.in/yaml.v2"

	"github.com/juju/juju/core/secrets"
	"github.com/juju/juju/secrets/provider"
)

var logger = loggo.GetLogger("juju.secrets.vault")

const (
	// Store is the name of the Kubernetes secrets store.
	Store = "vault"
)

// NewProvider returns a Vault secrets provider.
func NewProvider() provider.SecretStoreProvider {
	return vaultProvider{Store}
}

type vaultProvider struct {
	name string
}

func (p vaultProvider) Type() string {
	return p.name
}

// Initialise sets up a kv store mounted on the model uuid.
func (p vaultProvider) Initialise(m provider.Model) error {
	cfg, err := p.adminConfig(m)
	if err != nil {
		return errors.Trace(err)
	}
	client, err := p.newStore(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	sys := client.client.Sys()
	ctx := context.Background()

	mounts, err := sys.ListMountsWithContext(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	logger.Debugf("kv mounts: %v", mounts)
	modelUUID := cfg.Params["model-uuid"].(string)
	if _, ok := mounts[modelUUID]; !ok {
		err = sys.MountWithContext(ctx, modelUUID, &api.MountInput{
			Type:    "kv",
			Options: map[string]string{"version": "1"},
		})
		if !isAlreadyExists(err, "path is already in use") {
			return errors.Trace(err)
		}
	}
	return nil
}

// CleanupModel deletes all secrets and policies associated with the model.
func (p vaultProvider) CleanupModel(m provider.Model) error {
	cfg, err := p.adminConfig(m)
	if err != nil {
		return errors.Trace(err)
	}
	k, err := p.newStore(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	sys := k.client.Sys()

	// First remove any policies.
	ctx := context.Background()
	policies, err := sys.ListPoliciesWithContext(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	for _, p := range policies {
		if strings.HasPrefix(p, "model-"+k.modelUUID) {
			if err := sys.DeletePolicyWithContext(ctx, p); err != nil {
				if isNotFound(err) {
					continue
				}
				return errors.Annotatef(err, "deleting policy %q", p)
			}
		}
	}

	// Now remove any secrets.
	s, err := k.client.Logical().ListWithContext(ctx, k.modelUUID)
	if err != nil {
		return errors.Trace(err)
	}
	keys, ok := s.Data["keys"].([]interface{})
	if !ok {
		return nil
	}
	for _, id := range keys {
		err = k.client.KVv1(k.modelUUID).Delete(ctx, fmt.Sprintf("%s", id))
		if err != nil && !isNotFound(err) {
			return errors.Annotatef(err, "deleting secret %q", id)
		}
	}
	return nil
}

// CleanupSecrets removes policies associated with the removed secrets.
func (p vaultProvider) CleanupSecrets(m provider.Model, tag names.Tag, removed provider.SecretRevisions) error {
	cfg, err := p.adminConfig(m)
	if err != nil {
		return errors.Trace(err)
	}
	client, err := p.newStore(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	sys := client.client.Sys()

	isRelevantPolicy := func(p string) bool {
		for id := range removed {
			if strings.HasPrefix(p, fmt.Sprintf("model-%s-%s-", m.UUID(), id)) {
				return true
			}
		}
		return false
	}

	ctx := context.Background()
	policies, err := sys.ListPoliciesWithContext(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	for _, p := range policies {
		if isRelevantPolicy(p) {
			if err := sys.DeletePolicyWithContext(ctx, p); err != nil {
				if isNotFound(err) {
					continue
				}
				return errors.Annotatef(err, "deleting policy %q", p)
			}
		}
	}
	return nil
}

type vaultConfig struct {
	Endpoint      string   `yaml:"endpoint" json:"endpoint"`
	Namespace     string   `yaml:"namespace" json:"namespace"`
	Token         string   `yaml:"token" json:"token"`
	CACert        string   `yaml:"ca-cert" json:"ca-cert"`
	ClientCert    string   `yaml:"client-cert" json:"client-cert"`
	ClientKey     string   `yaml:"client-key" json:"client-key"`
	TLSServerName string   `yaml:"tls-server-name" json:"tls-server-name"`
	Keys          []string `yaml:"keys" json:"keys"`
}

// adminConfig returns the config needed to create a vault secrets store client
// with full admin rights.
func (p vaultProvider) adminConfig(m provider.Model) (*provider.StoreConfig, error) {
	cfg, err := m.Config()
	if err != nil {
		return nil, errors.Trace(err)
	}
	vaultCfgStr := cfg.SecretStoreConfig()
	if vaultCfgStr == "" {
		return nil, errors.NotValidf("empty vault config")
	}
	var vaultCfg vaultConfig
	if errJ := json.Unmarshal([]byte(vaultCfgStr), &vaultCfg); errJ != nil {
		if errY := yaml.Unmarshal([]byte(vaultCfgStr), &vaultCfg); errY != nil {
			return nil, errors.NewNotValid(errY, "invalid vault config")
		}
	}
	modelUUID := cfg.UUID()
	storeCfg := &provider.StoreConfig{
		StoreType: Store,
		Params: map[string]interface{}{
			"controller-uuid": m.ControllerUUID(),
			"model-uuid":      modelUUID,
			"endpoint":        vaultCfg.Endpoint,
			"namespace":       vaultCfg.Namespace,
			"token":           vaultCfg.Token,
			"ca-cert":         vaultCfg.CACert,
			"client-cert":     vaultCfg.ClientCert,
			"client-key":      vaultCfg.ClientKey,
			"tls-server-name": vaultCfg.TLSServerName,
		},
	}
	// If keys are provided, we need to unseal the vault.
	// (If not, the vault needs to be unsealed already).
	if len(vaultCfg.Keys) == 0 {
		return storeCfg, nil
	}

	vaultClient, err := p.newStore(storeCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	sys := vaultClient.client.Sys()
	for _, key := range vaultCfg.Keys {
		_, err := sys.Unseal(key)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return storeCfg, nil
}

// StoreConfig returns the config needed to create a vault secrets store client.
func (p vaultProvider) StoreConfig(m provider.Model, tag names.Tag, owned provider.SecretRevisions, read provider.SecretRevisions) (*provider.StoreConfig, error) {
	adminUser := tag == nil
	// Get an admin store client so we can set up the policies.
	storeCfg, err := p.adminConfig(m)
	if err != nil {
		return nil, errors.Trace(err)
	}
	store, err := p.newStore(storeCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	sys := store.client.Sys()

	ctx := context.Background()
	modelUUID := m.UUID()
	var policies []string
	if adminUser {
		// For admin users, all secrets for the model can be read.
		rule := fmt.Sprintf(`path "%s/*" {capabilities = ["read"]}`, modelUUID)
		policyName := fmt.Sprintf("model-%s-read", modelUUID)
		err = sys.PutPolicyWithContext(ctx, policyName, rule)
		if err != nil {
			return nil, errors.Annotatef(err, "creating read policy for model %q", modelUUID)
		}
		policies = append(policies, policyName)
	} else {
		// Agents can create new secrets in the model.
		rule := fmt.Sprintf(`path "%s/*" {capabilities = ["create"]}`, modelUUID)
		policyName := fmt.Sprintf("model-%s-create", modelUUID)
		err = sys.PutPolicyWithContext(ctx, policyName, rule)
		if err != nil {
			return nil, errors.Annotatef(err, "creating create policy for model %q", modelUUID)
		}
		policies = append(policies, policyName)
	}
	// Any secrets owned by the agent can be updated/deleted etc.
	logger.Debugf("owned secrets: %#v", owned)
	for id := range owned {
		rule := fmt.Sprintf(`path "%s/%s-*" {capabilities = ["create", "read", "update", "delete", "list"]}`, modelUUID, id)
		policyName := fmt.Sprintf("model-%s-%s-owner", modelUUID, id)
		err = sys.PutPolicyWithContext(ctx, policyName, rule)
		if err != nil {
			return nil, errors.Annotatef(err, "creating owner policy for %q", id)
		}
		policies = append(policies, policyName)
	}

	// Any secrets consumed by the agent can be read etc.
	logger.Debugf("consumed secrets: %#v", read)
	for id := range read {
		rule := fmt.Sprintf(`path "%s/%s-*" {capabilities = ["read"]}`, modelUUID, id)
		policyName := fmt.Sprintf("model-%s-%s-read", modelUUID, id)
		err = sys.PutPolicyWithContext(ctx, policyName, rule)
		if err != nil {
			return nil, errors.Annotatef(err, "creating read policy for %q", id)
		}
		policies = append(policies, policyName)
	}
	s, err := store.client.Auth().Token().Create(&api.TokenCreateRequest{
		TTL:             "10m", // 10 minutes for now, can configure later.
		NoDefaultPolicy: true,
		Policies:        policies,
	})
	if err != nil {
		return nil, errors.Annotate(err, "creating secret access token")
	}
	storeCfg.Params["token"] = s.Auth.ClientToken

	return storeCfg, nil
}

// NewVaultClient is patched for testing.
var NewVaultClient = vault.NewClient

// NewStore returns a vault backed secrets store client.
func (p vaultProvider) NewStore(cfg *provider.StoreConfig) (provider.SecretsStore, error) {
	return p.newStore(cfg)
}

func (p vaultProvider) newStore(cfg *provider.StoreConfig) (*vaultStore, error) {
	modelUUID := cfg.Params["model-uuid"].(string)
	address := cfg.Params["endpoint"].(string)

	var clientCertPath, clientKeyPath string
	clientCert, _ := cfg.Params["client-cert"].(string)
	clientKey, _ := cfg.Params["client-key"].(string)
	if clientCert != "" && clientKey == "" {
		return nil, errors.NotValidf("vault config missing client key")
	}
	if clientCert == "" && clientKey != "" {
		return nil, errors.NotValidf("vault config missing client certificate")
	}
	if clientCert != "" {
		clientCertFile, err := os.CreateTemp("", "client-cert")
		if err != nil {
			return nil, errors.Annotate(err, "creating client cert file")
		}
		defer func() { _ = clientCertFile.Close() }()
		clientCertPath = clientCertFile.Name()
		if _, err := clientCertFile.Write([]byte(clientCert)); err != nil {
			return nil, errors.Annotate(err, "writing client cert file")
		}

		clientKeyFile, err := os.CreateTemp("", "client-key")
		if err != nil {
			return nil, errors.Annotate(err, "creating client key file")
		}
		defer func() { _ = clientKeyFile.Close() }()
		clientKeyPath = clientKeyFile.Name()
		if _, err := clientKeyFile.Write([]byte(clientKey)); err != nil {
			return nil, errors.Annotate(err, "writing client key file")
		}
	}

	tlsConfig := vault.TLSConfig{
		TLSConfig: &api.TLSConfig{
			CACertBytes:   []byte(cfg.Params["ca-cert"].(string)),
			ClientCert:    clientCertPath,
			ClientKey:     clientKeyPath,
			TLSServerName: cfg.Params["tls-server-name"].(string),
		},
	}
	c, err := NewVaultClient(address,
		&tlsConfig,
		vault.WithAuthToken(cfg.Params["token"].(string)),
	)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if ns := cfg.Params["namespace"].(string); ns != "" {
		c.SetNamespace(ns)
	}
	return &vaultStore{modelUUID: modelUUID, client: c}, nil
}

type vaultStore struct {
	modelUUID string
	client    *vault.Client
}

// GetContent implements SecretsStore.
func (k vaultStore) GetContent(ctx context.Context, providerId string) (_ secrets.SecretValue, err error) {
	defer func() {
		err = maybePermissionDenied(err)
	}()

	s, err := k.client.KVv1(k.modelUUID).Get(ctx, providerId)
	if err != nil {
		return nil, errors.Annotatef(err, "getting secret %q", providerId)
	}
	val := make(map[string]string)
	for k, v := range s.Data {
		val[k] = fmt.Sprintf("%s", v)
	}
	return secrets.NewSecretValue(val), nil
}

// DeleteContent implements SecretsStore.
func (k vaultStore) DeleteContent(ctx context.Context, providerId string) (err error) {
	defer func() {
		err = maybePermissionDenied(err)
	}()

	err = k.client.KVv1(k.modelUUID).Delete(ctx, providerId)
	if isNotFound(err) {
		return nil
	}
	return err
}

// SaveContent implements SecretsStore.
func (k vaultStore) SaveContent(ctx context.Context, uri *secrets.URI, revision int, value secrets.SecretValue) (_ string, err error) {
	defer func() {
		err = maybePermissionDenied(err)
	}()

	path := uri.Name(revision)
	val := make(map[string]interface{})
	for k, v := range value.EncodedValues() {
		val[k] = v
	}
	err = k.client.KVv1(k.modelUUID).Put(ctx, path, val)
	return path, errors.Annotatef(err, "saving secret content for %q", uri)
}
