function threadSummarizeSystemPrompt(): string {
  return `You are tasked with summarizing conversation thread.` +
    `The goal is to extract the main points and present them in a concise, easy-to-read format.`;
}

export { threadSummarizeSystemPrompt };
