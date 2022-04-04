# Docs

We use the handy theme [just the docs](https://just-the-docs.github.io/just-the-docs/docs/navigation-structure/) for our more extensive documentation, and host it on github pages [here](https://smartcontractkit.github.io/integrations-framework/). Anything about the framework that can't be covered in our short-and-sweet [README](../README.md) should find a home in these docs.

## Local Development

1. Check if you have ruby installed, `ruby -v`. If not, [install it](https://www.ruby-lang.org/en/documentation/installation/).
2. `cd docs/` if you're not already there.
3. `bundle install`
4. `bundle exec jekyll serve --trace --livereload`
5. Visit `http://127.0.0.1:4000/` for a local version of the site

The local version should update anytime you save changes to the site files.

## Basic Style Guidelines

* Please use folders for top-level categories, even if they only include a single `index.md` file. This helps keep pages a bit more organized.
* `Chainlink` should be capitalized when possible, unless following a programming languages capitalization conventions in a code sample.
* Try using [Grammarly](https://app.grammarly.com/) or a similar service for spell-check and clarity suggestions.