"""Opal
Usage:
  opal (<fpath>)
  opal (-h | --help)

Options:
  -h --help     Show this screen.
"""

from docopt import docopt

import logging
from vault import ObsidianVault

logging.basicConfig(level=logging.INFO)

def main(dpath: str) -> None:
    vault = ObsidianVault(dpath)

    vault.validate_title()
    vault.fix_frontmatter()
    vault.fix_title()
    vault.flag_small_subgraphs()
    vault.validate_orphans()


if __name__ == '__main__':
    args = docopt(__doc__, version='Opal v0.1')
    main(args['<fpath>'])
