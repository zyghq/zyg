function threadSummarizeSystemPrompt(): string {
  return (
    `You are tasked with summarizing conversation thread.` +
    `The goal is to extract the main points and present them in a concise, easy-to-read format.`
  );
}

function threadCreateGitHubIssueSystemPrompt(): string {
  return (
    `You are an AI assistant specializing in creating GitHub issues from customer support chat transcripts. ` +
    `Your goal is to generate concise and actionable GitHub issues that engineers can use to ` +
    `resolve customer problems, address feature requests, and fix bugs efficiently.`
  );
}

export { threadSummarizeSystemPrompt, threadCreateGitHubIssueSystemPrompt };
