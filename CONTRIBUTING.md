# Contributing

We're so excited you're interested in helping with SuperTokens! We are happy to help you get started, even if you don't have any previous open-source experience :blush:

## New to Open Source?

1. Take a look at [How to Contribute to an Open Source Project on GitHub](https://egghead.io/courses/how-to-contribute-to-an-open-source-project-on-github)
2. Go through the [SuperTokens Code of Conduct](https://github.com/supertokens/supertokens-golang/blob/master/CODE_OF_CONDUCT.md)

## Where to ask Questions?

1. Check our [Github Issues](https://github.com/supertokens/supertokens-golang/issues) to see if someone has already answered your question.
2. Join our community on [Discord](https://supertokens.io/discord) and feel free to ask us your questions

## Development Setup

You will need to setup the `supertokens-core` in order to to run the `supertokens-golang` tests, you can setup `supertokens-core` by following this [guide](https://github.com/supertokens/supertokens-core/blob/master/CONTRIBUTING.md#development-setup)  
**Note: If you are not contributing to the `supertokens-core` you can skip steps 1 & 4 under Project Setup of the `supertokens-core` contributing guide.**

### Prerequisites

-   OS: Linux, macOS or [WSL](https://docs.microsoft.com/en-us/windows/wsl/about)
-   Go
-   IDE: [VSCode](https://code.visualstudio.com/download)(recommended) or equivalent IDE

### Project Setup

1. Fork the [supertokens-golang](https://github.com/supertokens/supertokens-golang) repository
2. Clone the forked repository in the parent directory of the previously setup `supertokens-root`.  
   `supertokens-golang` and `supertokens-root` should exist side by side within the same parent directory
3. `cd supertokens-golang`
4. You should have a go setup on your local machine

## Modifying Code

1. Open the `supertokens-golang` project in your IDE and you can start modifying the code

## Testing

1. Navigate to the `supertokens-root` repository
2. Start the testing environment  
   `./startTestingEnv --wait`
3. Navigate to the `supertokens-golang` repository  
   `cd ../supertokens-golang/`
4. Run all tests, [count=1 ensures tests are not cached]
   `go test  ./... -p 1 -v count=1`
5. If all tests pass the output should be:

## Pull Request

1. Before submitting a pull request make sure all tests have passed
2. Reference the relevant issue or pull request and give a clear description of changes/features added when submitting a pull request
3. Make sure the PR title follows [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification

## SuperTokens Community

SuperTokens is made possible by a passionate team and a strong community of developers. If you have any questions or would like to get more involved in the SuperTokens community you can check out:

-   [Github Issues](https://github.com/supertokens/supertokens-golang/issues)
-   [Discord](https://supertokens.io/discord)
-   [Twitter](https://twitter.com/supertokensio)
-   or [email us](mailto:team@supertokens.io)

Additional resources you might find useful:

-   [SuperTokens Docs](https://supertokens.io/docs/community/getting-started/installation)
-   [Blog Posts](https://supertokens.io/blog/)