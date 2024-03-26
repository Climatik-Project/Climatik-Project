# Define the JSON schema for the PowerCappingConfig CRD
POWER_CAPPING_CONFIG_SCHEMA = {
    "type": "object",
    "properties": {
        "powerCapLimit": {
            "type": "integer",
            "minimum": 0
        },
        "scaleObjectRefs": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "apiVersion": {
                        "type": "string"
                    },
                    "kind": {
                        "type": "string"
                    },
                    "metadata": {
                        "type": "object",
                        "properties": {
                            "name": {
                                "type": "string"
                            }
                        },
                        "required": ["name"]
                    }
                },
                "required": ["apiVersion", "kind", "metadata"]
            }
        }
    },
    "required": ["powerCapLimit", "scaleObjectRefs"]
}
