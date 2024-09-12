# GIt-Commit (gic)

Reducing cognitive load by automating commit message generation, allowing developers to focus on coding instead of crafting messages. A tool that helps developers generate git commit messages based on the `git diff` of staged files, following instructions. It's ideal for use alongside [Semantic Release](https://github.com/semantic-release/semantic-release).

## AzureAD

Remember to assign the `Azure Cognitive Openai User` role to any user that is going to consume the resource.

In your `.gic` config you would add:

```yaml
connection_type: "azure_ad"
azure_endpoint: "https://<endpoint>.openai.azure.com/"
model_deployment_name: "<model name>"
```

## Ollama Locally in your devcontainer (or any machine)

Here is an example to run ollama in your devcontainer and pulling phi3.5 image.

```json
...
"features": {
    "ghcr.io/prulloac/devcontainer-features/ollama:1": {
        "pull": "phi3.5"
    },
},
...
```

In your `.gic` config you would add:

```yaml
connection_type: "ollama"
azure_endpoint: "http://127.0.0.1:11434/"
model_deployment_name: "phi3.5"
```

>[!CAUTION]
>When choosing a model validate the instructions and generation matches what you expect. As an example, noticed phi3.5 didn't really generated a commit message with the requested instructions. More details in [here](#different-outputs-per-model).

## Different outputs per model

Instructions:

```yaml
llm_instructions: |

```

### phi3.5

```bash
```

### gpt-4o-mini

```bash
```
