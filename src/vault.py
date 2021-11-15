
from note import ObsidianNote
import re
import glob
from posixpath import basename
import logging


class ObsidianVault:
    """Represents an obsidian vault
    """
    dpath: str

    def __init__(self, dpath: str):
        self.dpath = dpath

    def list_markdown(self) -> list[str]:
        """List markdown files in a directory; do not consider deeper markdown since I only use one folder.
        """
        return [basename(file)for file in glob.glob(self.dpath + "/*.md")]

    def validate_title(self):
        """Validate the obsidian notes title is in <numeric-time> - <title> format
        """
        fnames = self.list_markdown()
        misnamed = {fname for fname in fnames if not re.match(
            '[0-9]{12} - ', fname)}

        if len(misnamed) > 0:
            logging.error(
                f'💎 there were {len(misnamed)} misnamed files: \n' + '\n'.join(misnamed))
        else:
            logging.info('💎 files named correctly')

    def fix_frontmatter(self):
        """Ensure all files have frontmatter present, including an alias and tags.
        """
        fnames = self.list_markdown()

        for fname in fnames:
            note = ObsidianNote(self.dpath, fname)
            note.fix_frontmatter()

    def fix_title(self):
        """Fix the note title"""
        fnames = self.list_markdown()

        for fname in fnames:
            note = ObsidianNote(self.dpath, fname)
            note.fix_title()

    def validate_orphans(self):
        """Validate there are no orphan notes present."""
        fnames = self.list_markdown()
        orphans = 0

        for fname in fnames:
            note = ObsidianNote(self.dpath, fname)
            count = note.wikilink_count()

            if count == 0:
                orphans += 1
                print(fname)

        if orphans > 0:
            raise Exception(f'there were {orphans} orphans')

    def flag_small_subgraphs(self):
        fnames = self.list_markdown()

        for fname in fnames:
            note = ObsidianNote(self.dpath, fname)
