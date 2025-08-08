# Token Optimization Strategies

When working with LLM agents, minimizing token usage is crucial for efficiency and cost-effectiveness. The `ffs` library is designed with token optimization in mind. Here are some strategies you can employ:

## 1. Use Diffs for Updates

Instead of sending the entire file content to an LLM for modification, you can send a diff of the changes. This significantly reduces the number of tokens required.

The `diff` package in `ffs` can be used to generate and apply patches, ensuring that only the changed lines are processed.

## 2. Chunk Large Files

For very large files, consider breaking them down into smaller chunks before sending them to the LLM. This can help you stay within the token limits of the model and can also improve the quality of the suggestions.

The `core` package can be used to read and write specific parts of a file, allowing you to implement a chunking strategy.

## 3. Minimize Context

When interacting with an LLM, provide only the necessary context for the task at hand. Avoid sending irrelevant parts of the codebase or documentation.

By using the `ffs` library to selectively read files and directories, you can build a focused context for the LLM, which will result in more accurate and efficient suggestions.
