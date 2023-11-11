from lib.base import BaseTest
from lib.modules.tests import ModulesTest
from lib.providers.tests import ProvidersTest
import pkgutil
__path__ = pkgutil.extend_path(__path__, __name__)
__all__ = ['BaseTest', 'ModulesTest', 'ProvidersTest']
