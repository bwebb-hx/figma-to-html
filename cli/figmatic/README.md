# Figmatic: A CLI for Generating Code from Figma Designs

This CLI is for facilitating the creation of code from Figma designs. As of now, it only generates HTML and CSS code, but perhaps in the future the scope can be expanded to other programming languages and front-end frameworks, such as React.

## Requirements

You must have `claude` (Claude Code) installed and findable in your system's PATH variable.

## Setup

Install this CLI globally:

```bash
go install github.com/bwebb-hx/figma-to-html/cli/figmatic
```

### Required! MCP Server Setup

You must have the Figma MCP server registered in Claude Code, under the **USER** scope:

```bash
claude mcp add figma-mcp-1 -s user -- npx -y figma-developer-mcp --figma-api-key=$FIGMA_ACCESS_TOKEN --stdio
```

> For now, we have a hard-coded MCP name of "figma-mcp-1". However, in the future I may add config for customizing this name.

If you are unsure if your MCP server is set up right, try entering the following command:

```bash
claude mcp list
```

If it's correctly set under the user scope, it should show your MCP config regardless of what your current directory is in the terminal.

## Usage

For detailed usage instructions, enter `figmatic -h` in the terminal.

But, here's an example of generating a figma design into code:

```bash
figmatic gen --url "https://www.figma.com/design/AxOg233wgtoO27RFd8QIOv/M55_コーディングテスト用?node-id=1-1062" -t "$FIGMA_ACCESS_TOKEN" --sub-nodes
```

This command does the following:

```
gen:          command for generating code from figma designs
--url:        specify the URL of the figma design
-t:           your Figma access token
--sub-nodes:  tells the tool to generate the sub-nodes (one level below) first, and then combine those results. For improved accuracy.
```