# Using Claude Code to Generate HTML from Figma

## Adding an MCP Server to Claude Code

Related docs: https://docs.anthropic.com/en/docs/claude-code/mcp#configure-mcp-servers

Here's the basic command syntax for adding a new mcp server to claude code:

```
claude mcp add <name> -- <command> [args...]
```

> Note: in the documentation, they don't show the "--" in the "Basic Syntax" section, but they do in the example.
> I tried without the "--" first and it errored, so it seems that is probably required.

For our Figma MCP server command, we can enter something like the following:

```
claude mcp add figma-mcp-1 -- npx -y figma-developer-mcp --figma-api-key=FIGMA_ACCESS_TOKEN_HERE --stdio
```

> Make sure to sub in your actual figma access token

See the full documentation for other relevant MCP commands.

## Attempting to use Figma MCP Server with Claude Code

It looks like the MCP server is working, so I gave Claude Code a try by giving it a figma design link.

I've gotten an error once the MCP server started downloading the figma design data:

```
⏺ figma-mcp-1:get_figma_data (MCP)(fileKey: "AxOg233wgtoO27RFd8QIOv", nodeId: "1-1061")
  ⎿  Error: MCP tool "get_figma_data" response (27548 tokens) exceeds maximum allowed tokens (25000). Please use
     pagination, filtering, or limit parameters to reduce the response size.
```

I simply ran the same thing again (I think? maybe claude code adjusted the approach this time?) and it started working now.

It successfully generated some HTML, but it seems to have only created a small amount of the design - mainly the top section and only a tiny bit of it.

I asked it to take another look at the design and try again. 

> Did it fail to download all of the design due to that previous error about "maximum allowed tokens"?

After it finished the second iteration, much more of the design was now there. Awesome! But, still a significant amount of inaccuracies. Still, better than nothing. I wonder if I'll need to do several iterations to slowly get closer and closer to the correct design? Let's try one more iteration.

The third iteration has taken significantly longer (around 300s), and I did notice an error message in the output:

```
⏺ Update(index.html)
  ⎿  Error: Found 4 matches of the string to replace, but replace_all is false. To replace all occurrences, set
     replace_all to true. To replace only one occurrence, please provide more context to uniquely identify the
     instance.
     String:                         <div class="news-item">
                                 <div class="news-date">2025年6月9日</div>
                                 <div class="news-content">
                                     <span class="news-label">NEWS</span>
                                     <div class="news-text">Consolidated Financial Results for the Fiscal Year Ended
     March 31, 2025 (Under Japanese GAAP)</div>
                                 </div>
                             </div>
```

It seemed like Claude Code was able to move on from this and try to resolve the issue. At some point it did finish up, but not clear what the hold up was.
Upon reviewing the 3rd iteration, it did seem to fix a couple things and get a bit closer to the figma design, but I think the improvements have slowed down.

Maybe this tool will just be mainly useful for getting designs started, and then manual intervention will be necessary? Or maybe it's better to build smaller pieces bit by bit, and put it all together?

After this session with 3 iterations of asking claude code to generate HTML for the figma design, it looks like in total we incurred about $1.50 of charges:

```
Total cost:            $1.49
Total duration (API):  21m 25.5s
Total duration (wall): 37m 7.4s
Total code changes:    781 lines added, 57 lines removed
Usage by model:
    claude-3-5-haiku:  68.9k input, 2.5k output, 0 cache read, 0 cache write
       claude-sonnet:  198 input, 23.6k output, 1.2m cache read, 186.9k cache write
```

## Creating a Script to Automate Figma Code Generation

The main benefit of claude code over cursor is that we can potentially use this CLI version of claude to create scripts and automate things. Let's explore the possibilities here.

First, let's think of how scripting or automation would be useful in terms of converting figma designs to code:

1. Given a list of URLs, generate the HTML/CSS for each design.
- this would be especially useful if there were a lot of designs we need to generate, and each design might take some time.
- give the script a list of URLs and let it run (with auto-accept edits turned on). come back later and find all the designs converted into HTML!
- if generating smaller sub-sections of the design is more accurate, we could also divide a page of a figma design into sub-sections, get the URLs for each, have those generated individually, and then put them all together after.

2. An iterative claude code session?
- I observed in my first test that the first iteration was missing a large amount of content and generally inaccurate.
- So, perhaps we could make a script that causes claude to iterate a few times on the same figma design, gradually improving it without needing user intervention.
- We could give claude a prompt each time, perhaps asking it to "look for inaccuracies and fix them". And set a number of iterations to repeat this?

> in order to use the figma MCP server in claude code while NOT in interactive mode, you need to give claude code permission.
> 