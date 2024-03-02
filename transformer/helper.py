from pre_install import encoding

def get_token_length(test: str) -> int:
    return len(encoding.encode(test))