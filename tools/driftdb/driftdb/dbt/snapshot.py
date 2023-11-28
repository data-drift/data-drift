import json


def get_snapshot_nodes() -> list[str]:
    project_dir = "."

    with open(f"{project_dir}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

        snapshot_nodes = [node for node in manifest["nodes"].values() if (node["resource_type"] == "snapshot")]
        return [node["unique_id"] for node in snapshot_nodes]
