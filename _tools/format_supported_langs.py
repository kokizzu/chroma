#!/usr/bin/env python3

import argparse
import re
import sys
from collections import defaultdict

BEGIN_MARKER = "<!-- supported-languages:begin -->"
END_MARKER = "<!-- supported-languages:end -->"


def parse_lexers(lines: list[str]) -> list[str]:
	"""Parse the output of chroma --list and return a list of lexer names"""

	lexer_name_re: re.Pattern[str] = re.compile(r"^  ([^:\s].*?)\s*$")
	lexers: list[str] = []
	in_lexers = False

	for line in lines:
		line = line.rstrip()
		if line.startswith("lexers:"):
			in_lexers = True
			continue
		if not in_lexers:
			continue

		# stop when we hit styles/formatters/etc
		if line.startswith("styles:") or line.startswith("formatters:"):
			break

		match: re.Match[str] | None = lexer_name_re.match(line)
		if match:
			name: str | None = match.group(1)
			if name:
				lexers.append(name)
	return lexers


def group_by_prefix(lexers: list[str]) -> dict[str, list[str]]:
	"""Given a list of lexer names, return a dictionary mapping prefixes
	to lists of lexers that begin with that prefix"""
	groups: defaultdict[str, list[str]] = defaultdict(list[str])
	for name in lexers:
		prefix: str = name[0].upper()
		groups[prefix].append(name)
	# sort alphabetically
	for k in groups:
		groups[k] = sorted(groups[k], key=lambda s: s.lower())
	return dict(sorted(groups.items()))


def emit_markdown(groups: dict[str, list[str]]) -> str:
	lines: list[str] = []
	longest = 0
	for prefix, lexers in groups.items():
		joined: str = ", ".join(lexers)
		l: int = len(joined)
		if l > longest:
			longest: int = l
		lines.append(f"|   {prefix}    | {joined}")
	splitter = f"| :----: | {longest * '-'}"
	markdown: list[str] = ["| Prefix | Language", splitter]
	markdown.extend(lines)
	return "\n".join(markdown)


def update_file(path: str, markdown: str) -> None:
	"""Replace the table between the supported-languages markers in path"""
	with open(path) as f:
		content: str = f.read()
	try:
		begin: int = content.index(BEGIN_MARKER) + len(BEGIN_MARKER)
		end: int = content.index(END_MARKER)
	except ValueError:
		sys.exit(f"{path}: missing {BEGIN_MARKER} / {END_MARKER} markers")
	with open(path, "w") as f:
		f.write(content[:begin] + "\n" + markdown + "\n" + content[end:])


if __name__ == "__main__":
	parser = argparse.ArgumentParser(
		description="Parse chroma --list piped from stdin and emit a markdown table for the README"
	)
	parser.add_argument(
		"--update",
		metavar="FILE",
		help="replace the table between the supported-languages markers in FILE instead of printing it",
	)
	args = parser.parse_args()

	if sys.stdin.isatty():
		parser.print_help()
		print("\nRecommended usage (from repo root):")
		print(
			"go -C cmd/chroma run . --list | _tools/format_supported_langs.py --update README.md"
		)
		exit(1)

	lines: list[str] | None = sys.stdin.readlines()
	if not lines:
		exit(1)
	lexers: list[str] = parse_lexers(lines)
	groups: dict[str, list[str]] = group_by_prefix(lexers)
	markdown: str = emit_markdown(groups)
	if args.update:
		update_file(args.update, markdown)
	else:
		print(markdown)
