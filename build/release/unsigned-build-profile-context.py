"""Fixed unsigned-build source binding using the reviewed same-user observer."""
import importlib.util
from pathlib import Path
import sys
spec = importlib.util.spec_from_file_location('pscan_native_context', Path(__file__).with_name('native-profile-context.py'))
context = importlib.util.module_from_spec(spec)
spec.loader.exec_module(context)
context.SOURCES = (
    'unsigned-build-profile-context.py', 'unsigned-build-profile-driver.ps1',
    'unsigned-build-stage.ps1', 'native-profile-context.py', 'native-host-facts.py',
    'docker-execution.ps1', 'execution-profile.ps1', 'docker-container-lifecycle.ps1', 'invoke-docker-boundary.ps1', 'native-fixture-diagnostics.ps1',
    'test-docker-execution.ps1', 'test-linux-handoff.ps1', 'test-image-admission.ps1', 'image-admission.ps1',
    'test-cache-boundary.ps1', 'cache-canary.ps1', 'host-cache-canary.ps1',
    'admit-image.ps1', 'acquire.ps1', 'build.ps1', 'build-validation.ps1', 'invoke-exact-build.ps1',
    'source-trust.ps1', 'test-crlf-shell-payloads.ps1', 'compare-builds.ps1')
if __name__ == '__main__':
    sys.exit(context.main())
