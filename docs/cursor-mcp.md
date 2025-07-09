# Using Cursor MCP to Generate Figma Designs

This guide is for setting up an MCP server with Cursor, and enabling Cursor to create the code for Figma designs.

I've started by basing my approach on this Youtube video: https://www.youtube.com/watch?v=X-aX1TuGP0s

## Setting up MCP Server

I've found this MCP server for Figma on Github: https://github.com/GLips/Figma-Context-MCP

To add an MCP server in Cursor, go to Cursor Settings > Tools & Integrations > Add MCP Server.
Alternatively, you can create a `/.cursor` directory in your project root directory, and then inside create a `mcp.json` file.

Below is an example of what `mcp.json` might look like:

```json
{
    "mcpServers": {
        "Framelink Figma MCP": {
            "command": "npx",
            "args": ["-y", "figma-developer-mcp", "--figma-api-key=FIGMA_API_KEY_HERE", "--stdio"]
        }
    }
}
```

### Personal Access Tokens

This page explains how to generate access tokens: https://help.figma.com/hc/en-us/articles/8085703771159-Manage-personal-access-tokens

For some reason figma has some bugs and will fail to generate tokens sometimes. It doesn't say why, but I guess it has something to do with the specific permissions granted when generating a token.

I successfully generated a token with the following permissions:

- code connect: write
- dev resources: read only
- file content: read only
- file metadata: read only
- library assets: read only
- library content: read only

