import os
import runpy
import sys

import dbt_snapshot_analysis


def main() -> None:
    streamlit_script_path = os.path.join(
        os.path.dirname(dbt_snapshot_analysis.__file__), "dbt_snapshot_analysis.py"
    )
    sys.argv = ["streamlit", "run", streamlit_script_path]
    runpy.run_module("streamlit", run_name="__main__")


if __name__ == "__main__":
    main()
