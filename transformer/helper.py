from pre_install import encoding, tokenizer, model, t5_model
import re


def remove_tags(text):
    """
    Remove HTML tags from a string.

    Args:
      text: The string to remove HTML tags from.

    Returns:
      The string with HTML tags removed.
    """
    return re.sub(r"<[^>]*?>", "", text)


def get_token_length(test: str) -> int:
    return len(encoding.encode(test))


def get_embedding(text: str) -> list:
    return list(model.encode(text))


content = """
    You have a vector database containing summarized vectors of GitHub repositories that users have starred.
    Now, a user has presented a question.
    Please help me refine the following description for a more accurate expression, reply in English:
    {question}
    """


def get_formatted_text(text: str) -> str:

    input_ids = tokenizer(content.format(question=text), return_tensors="pt").input_ids
    outputs = t5_model.generate(input_ids, max_new_tokens=60)

    return remove_tags(tokenizer.decode(outputs[0]))
