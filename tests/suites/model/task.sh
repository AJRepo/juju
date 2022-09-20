test_model() {
	if [ "$(skip 'test_model')" ]; then
		echo "==> TEST SKIPPED: model tests"
		return
	fi

	set_verbosity

	echo "==> Checking for dependencies"
	check_dependencies juju

	file="${TEST_DIR}/test-models.log"

	bootstrap "test-models" "${file}"

	# Tests that need to be run are added here.
	test_model_config
	test_model_migration
	test_model_migration_version
	test_model_migration_saas_common
	test_model_migration_saas_external
	test_model_multi
	test_model_metrics
	test_model_destroy

	destroy_controller "test-models"
}
