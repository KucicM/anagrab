import argparse
import string
import pathlib

ASCII_LETTERS = set(string.ascii_letters)
END_CHAR = "!"


def build_word_list(input_path: str) -> str:
    words = []
    print("load words")
    with open(input_path, "r") as f:
        words = f.readlines()

    print("building word list")
    words = filter(is_ascii_word, words)
    words = map(lambda x: x.lower(), words)
    words = map(lambda x: x.strip(), words)
    words = sorted(set(words), key=lambda w: sorted(w))

    return "\n".join(words)


def is_ascii_word(word: str) -> bool:
    return all(w in ASCII_LETTERS for w in word.strip())


def to_bin(bins: bytes):
    return repr(bins)[2:-1]


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("input_file_path")
    parser.add_argument("output_file_path")

    args = parser.parse_args()
    words = build_word_list(args.input_file_path)

    out_path = pathlib.Path(args.output_file_path)
    with open(out_path / "word_list.txt", "w") as f:
        f.write(words)
