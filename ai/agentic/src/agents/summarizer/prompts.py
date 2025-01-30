from typing import Optional, Dict, Callable
from pydantic import BaseModel

# Prompt version constants
EXECUTIVE_PROMPT = 1
ENTITY_PROMPT = 2
QA_PROMPT = 3
BULLET_POINT_PROMPT = 4


# TypedDict definitions for JSON output (with added strictness)
class ExecutiveSummary(BaseModel):
    title: str
    audience: str
    summary: str


class Entity(BaseModel):
    title: str
    entities: list[str]
    summary: str


class QAPair(BaseModel):
    question: str
    answer: str


class QuestionAnswerSummary(BaseModel):
    title: str
    qa_pairs: list[QAPair]
    summary: str


class BulletPointSummary(BaseModel):
    title: str
    bullet_points: list[str]
    summary: str


def create_executive_prompt(
    text: str, max_words: int, context: Optional[str] = None
) -> str:
    prompt = f"""# EXECUTIVE SUMMARY GENERATION
        ## Instructions
        You are an expert summarizer. Your task is to create a concise, well-structured executive summary of the following text in {max_words} words or less. Identify the intended audience and tailor the summary accordingly, prioritizing the most critical information within the word limit.

        ## Input Text
        ```
        {text}
        ```
        ## Context Information
        {f"Additional Context: {context}" if context else "No additional context provided."}

        ## Requirements
        1. Identify the main subject/topic and intended audience.
        2. Extract and prioritize the most crucial information.
        3. Maintain a professional tone appropriate for the identified audience.
        4. Ensure accuracy, completeness, and logical structure.
        5. Include a concise title that encapsulates the main topic.

        **Important:** Your response must be a valid JSON object. Do not include any introductory text, concluding phrases, or additional explanations outside of the JSON object.
        """
    return prompt


def create_entity_prompt(
    text: str, max_words: int, context: Optional[str] = None
) -> str:
    prompt = f"""# ENTITY-FOCUSED SUMMARY GENERATION
        ## Instructions
        As an expert summarizer, create an entity-dense summary of the following content in {max_words} words or less. The summary should seamlessly integrate key entities into a cohesive narrative.

        ## Input Text
        ```
        {text}
        ```
        ## Context Information
        {f"Additional Context: {context}" if context else "No additional context provided."}

        ## Requirements
        1. Identify approximately 5-10 key Descriptive Entities that are:
        - Relevant to the main content.
        - Specific (5 words or fewer).
        - Faithful to the source.

        2. Create a summary that:
        - Naturally weaves all identified entities into the narrative.
        - Is {max_words} words or fewer.
        - Is self-contained and easily understood.
        - Directly addresses the content without introductory phrases or self-reference.

        3. Improve density through:
        - Entity fusion and compression.
        - Removal of filler phrases.
        - Information-rich language.

        **Important:** Your response must be a valid JSON object. Do not include any introductory text, concluding phrases, or additional explanations outside of the JSON object.
        """
    return prompt


def create_question_answer_prompt(
    text: str, max_words: int, context: Optional[str] = None
) -> str:
    prompt = f"""# QUESTION & ANSWER EXTRACTION AND SUMMARY
        ## Instructions
        You are an expert in identifying and summarizing key information in a Question and Answer format. Your task is to extract the most important questions and answers from the following text and then create a concise summary based on these Q&As in {max_words} words or less.

        ## Input Text
        ```
        {text}
        ```
        ## Context Information
        {f"Additional Context: {context}" if context else "No additional context provided."}

        ## Requirements
        1. Identify 3-5 of the most important questions addressed in the text.
        2. Extract or infer the corresponding answers to these questions directly from the text.
        3. Ensure the answers are concise, accurate, and directly responsive to the questions.
        4. Create a summary that synthesizes the information from the Q&A pairs.
        5. The summary should be {max_words} words or fewer.

        **Important:** Your response must be a valid JSON object. Do not include any introductory text, concluding phrases, or additional explanations outside of the JSON object.
        """
    return prompt


def create_bullet_point_prompt(
    text: str, max_words: int, context: Optional[str] = None
) -> str:
    prompt = f"""# BULLET POINT SUMMARY GENERATION
        ## Instructions
        You are an expert summarizer skilled in creating concise bullet point summaries. Condense the following text into its most important points, using bullet points.

        ## Input Text
        ```
        {text}
        ```
        ## Context Information
        {f"Additional Context: {context}" if context else "No additional context provided."}

        ## Requirements
        1. Identify and extract the 5-7 most crucial points from the text.
        2. Present each point as a separate bullet point.
        3. Ensure each bullet point is concise, clear, and directly relevant to the main topic.
        4. Order the bullet points logically to reflect the flow or importance of the information.
        5. Each bullet point should be a single, succinct statement.
        6. Provide a title that reflects the subject of the text.
        7. The summary should be {max_words} words or fewer.

        **Important:** Your response must be a valid JSON object. Do not include any introductory text, concluding phrases, or additional explanations outside of the JSON object.
        """
    return prompt


# Dictionary mapping prompt versions to their respective functions
PROMPT_FUNCTIONS: Dict[int, Callable] = {
    EXECUTIVE_PROMPT: create_executive_prompt,
    ENTITY_PROMPT: create_entity_prompt,
    QA_PROMPT: create_question_answer_prompt,
    BULLET_POINT_PROMPT: create_bullet_point_prompt,
}
