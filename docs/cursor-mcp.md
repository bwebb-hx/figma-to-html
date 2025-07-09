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

If everything is set up correctly, you should see a row for the new MCP server in Cursor Settings now, and it should have a green light indicating connection is working.

### Figma Personal Access Token

This page explains how to generate access tokens: https://help.figma.com/hc/en-us/articles/8085703771159-Manage-personal-access-tokens

> For some reason figma was having some errors when I was trying to create access tokens at first. I assumed it might be due to the permissions I was choosing, but in the end things worked out fine. Not sure what the issue was.

I successfully generated a token with the following permissions:

- code connect: write
- dev resources: read only
- file content: read only
- file metadata: read only
- library assets: read only
- library content: read only

The above permissions seemed to work well enough. It's unclear to me though if I should get other permissions too.

## Generating Figma -> HTML with Cursor and MCP

At this point everything should be set up correctly, and we are ready to go.

All you have to do is simply paste a link to the figma design (or component, sub-section of the design, etc) in the Cursor Agent chat, and then tell it how you want the design implemented. 

### Test Run 1

I did this for my first test run:

```
Hi, can you create the linked Figma design in HTML/CSS/Javascript?  Don't use any SPA frameworks (e.g. React), just use plain HTML and CSS.

@https://www.figma.com/design/AxOg233wgtoO27RFd8QIOv/M55_%E3%82%B3%E3%83%BC%E3%83%86%E3%82%99%E3%82%A3%E3%83%B3%E3%82%AF%E3%82%99%E3%83%86%E3%82%B9%E3%83%88%E7%94%A8?node-id=1-1061&m=dev 
```

And, it worked pretty well!  It generated a relatively accurate version of the figma design, in only HTML and CSS. However, it did not include any Javascript, so things like button presses and pagination didn't work. Looking at the figma design though, I don't think any pagination was created, so I guess Cursor didn't have any pages to create besides the single one.  This might be solvable pretty easily though, once we tell Cursor to implement the Javascript code for navigation between pages and stuff like that.

Another issue was that none of the SVG images got downloaded. Cursor's AI Agent seems to have access to a `download_figma_images` function from the Figma MCP, but I guess it's not working properly with this design. According to the AI, the request is responding with "Success, 0 images downloaded".

I've tried a few things, like adding an Export option to the SVG layers in Figma, but it's still not working.

#### TODO

1. See if we can get SVG/other images to download with Cursor's AI Agent and the Figma MCP server.

2. Experiment with generating an entire figma design (all the pages) rather than just a single piece of the design.