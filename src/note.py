
from dataclasses import dataclass
import os
import re
from typing import Optional

import yaml


@dataclass
class ObsidianNote:
    fpath: str
    dpath: str
    name: str
    date: int

    def __init__(self, dpath: str, fpath: str):
        self.fpath = fpath
        self.dpath = dpath

        parts = fpath.split(' - ', maxsplit=1)

        if len(parts) != 2:
            raise Exception(f'misformatted name {fpath}')

        date, name = parts

        self.name = name
        self.date = int(date)

    def wikilink_count(self) -> int:
        """Count how many wikilinks are present in a document. Should probably switch to MD parsing..."""
        with open(os.path.join(self.dpath, self.fpath), 'r') as conn:
            content = conn.read()
            lines = content.split('\n')

            # -- not perfect...
            links = [re.match(r'.*\[\[.+\]\].*', line) for line in lines]
            links = [link for link in links if link is not None]

            return len(links)

    def find_title(self, content: str, fpath: str) -> dict:
        """Find the document's title tag"""
        h1 = [line for line in content.split('\n') if line.startswith('# ')]

        if len(h1) == 0:
            name = fpath.split(' - ', 1)[1]
            return {
                'title': re.sub(r'.md$', '', name),
                'hasTag': False
            }

        return {
            'title': h1[0][2:],
            'hasTag': True
        }

    def find_tags(self, content: str) -> list[str]:
        """Find tags within the document"""

        tags_regexp = r'tags: #.+'
        tags = [line for line in content.split(
            '\n') if re.match(tags_regexp, line)]

        if len(tags) > 0:
            for tagset in tags:
                for tag in tagset.replace('tags: ', '').split(', '):
                    yield tag.strip()

    def read_frontmatter(self, content: str):
        """Read YAML frontmatter from notes"""
        lines = [line for line in content.split('\n') if len(line.strip()) > 0]

        if len(lines) == 0:
            return None, None

        if not lines[0].startswith('---'):
            return None, None
        else:
            opened = False
            frontmatter = []

            line_idx = 0

            for line in lines:
                if line.startswith('---'):
                    if opened:
                        return yaml.load('\n'.join(frontmatter), Loader=yaml.Loader), line_idx
                    else:
                        opened = True
                else:
                    frontmatter.append(line)
                line_idx += 1

            # -- probably not frontmatter, but a divider
            return None, None

    def update_frontmatter(self, fm: Optional[dict], data: dict):
        """Merge existing frontmatter, if present, into the document-extracted frontmatter
        """
        updated = fm if fm else {}

        if not 'tags' in updated or updated['tags'] == None:
            updated['tags'] = []

        if not 'aliases' in updated or updated['aliases'] == None:
            updated['aliases'] = []

        tags = [updated['tags']] if isinstance(
            updated['tags'], str) else updated['tags']
        aliases = [updated['aliases']] if isinstance(
            updated['aliases'], str) else updated['aliases']

        updated['tags'] = list(set(tags + data['tags']))
        updated['aliases'] = list({alias.lower()
                                  for alias in [data['title']] + aliases})

        return updated

    def write_frontmatter(self, frontmatter: dict, content: str, end: Optional[int]) -> None:
        """Write frontmatter to a document, being mindful of existing frontmatter"""
        lines = content.split('\n')

        frontmatter_lines = ['---'] + \
            yaml.dump(frontmatter).split('\n') + ['---']

        frontmatter_lines = [line for line in frontmatter_lines if len(line.strip()) > 0]

        new_lines = frontmatter_lines + lines[end + 1:] if end else frontmatter_lines + lines

        with open(os.path.join(self.dpath, self.fpath), 'w') as conn:
            conn.write('\n'.join(new_lines))

    def fix_frontmatter(self):
        """Update or set a note's frontmatter.
        """
        with open(os.path.join(self.dpath, self.fpath), 'r') as conn:
            content = conn.read()
            fm, end = self.read_frontmatter(content)

            tags = list(self.find_tags(content))
            title = self.find_title(content, self.fpath)

            if fm:
                if end and end > 7:
                    raise Exception(f'huh? {self.fpath}')

            updated = self.update_frontmatter(fm, {
                'tags': tags,
                'title': title['title']
            })

        self.write_frontmatter(updated, content, end)

    def write_title(self, title: str, content: str, end: Optional[int]):
        """Write title to the document"""

        lines = content.split('\n')

        title_lines = [f'# {title}', '---']

        new_content = title_lines + \
            lines if not end else lines[0:end+1] +title_lines + lines[end+1:]

        with open(os.path.join(self.dpath, self.fpath), 'w') as conn:
            conn.write('\n'.join(new_content))

    def fix_title(self):
        """If a title is missing, add it based on the file-name"""
        with open(os.path.join(self.dpath, self.fpath), 'r') as conn:
            content = conn.read()
            title = self.find_title(content, self.fpath)

            _, end = self.read_frontmatter(content)

        if not title['hasTag']:
            self.write_title(title['title'], content, end)
