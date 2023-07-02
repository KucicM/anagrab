import argparse
import string
import heapq
import pathlib
from collections import Counter
from typing import Dict, Tuple

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


def build_compression_tree(words: str) -> Tuple[Dict[str, bytes], Dict[bytes, str]]:
    print("building compression tree")
    nodes = [(f, i, w) for i, (w, f) in enumerate(Counter(words).items())]
    nodes.append((1, -1, END_CHAR))
    heapq.heapify(nodes)

    while len(nodes) > 1:
        f1, i, n1 = heapq.heappop(nodes)
        f2, i, n2 = heapq.heappop(nodes)
        heapq.heappush(nodes, (f1 + f2, i, (n1, n2)))

    def build(node, tree: Dict[str, bytes], b=b""):
        if isinstance(node, str):
            tree[node] = b
            return
        build(node[0], tree, b + b"0")
        build(node[1], tree, b + b"1")

    tree: Dict[str, bytes] = dict()
    build(nodes[0][2], tree)
    return tree, {v: k for k, v in tree.items()}


def compress_words(word: str, tree: Dict[str, bytes]) -> bytes:
    print("compressing")
    encoded = [tree[c] for c in word + f"\n{END_CHAR}"]
    return b"".join(encoded)


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

    encode_tree, decode_tree = build_compression_tree(words)
    with open(out_path / "compression_tree.txt", "w") as f:
        parts = "|".join(["|".join([to_bin(k), v]) for k, v in decode_tree.items()])
        f.write(parts)

    e = compress_words(words, encode_tree)
    with open(out_path / "word_list.bin", "wb") as f:
        for i in range(0, len(e), 8):
            byte = e[i:i+8] + b"0" * (8 - len(e[i:i+8]))
            f.write(int(byte, 2).to_bytes(1, "big"))
