import importlib.util
import os
import sys


def import_user_defined_function(file_path, function_name):
    file_path = os.path.abspath(file_path)

    module_name = os.path.splitext(os.path.basename(file_path))[0]

    print(f"Importing {module_name} from {file_path}")

    spec = importlib.util.spec_from_file_location(module_name, file_path)
    module = importlib.util.module_from_spec(spec)
    sys.modules[module_name] = module
    spec.loader.exec_module(module)

    user_function = getattr(module, function_name, None)
    if not user_function:
        raise ImportError(f"Function {function_name} not found in {file_path}")

    return user_function
