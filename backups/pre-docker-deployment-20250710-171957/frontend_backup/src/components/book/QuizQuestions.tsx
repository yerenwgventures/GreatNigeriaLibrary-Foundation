import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { fetchQuizQuestions } from '../../features/books/booksSlice';

const QuizContainer = styled.div`
  margin: 2rem 0;
  padding: 1.5rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #16213e;
`;

const Title = styled.h3`
  font-size: 1.3rem;
  margin-bottom: 1rem;
  color: #16213e;
`;

const QuestionsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const QuestionItem = styled.div`
  background-color: white;
  border-radius: 4px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
`;

const QuestionText = styled.h4`
  font-size: 1.1rem;
  margin-bottom: 1rem;
  color: #16213e;
`;

const OptionsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
`;

const OptionItem = styled.label<{ isSelected: boolean; isCorrect?: boolean; isIncorrect?: boolean }>`
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  background-color: ${props => 
    props.isCorrect ? '#e6f7e6' : 
    props.isIncorrect ? '#ffebee' : 
    props.isSelected ? '#e3f2fd' : 
    '#f5f5f5'
  };
  border: 1px solid ${props => 
    props.isCorrect ? '#81c784' : 
    props.isIncorrect ? '#ef9a9a' : 
    props.isSelected ? '#90caf9' : 
    '#ddd'
  };
  
  &:hover {
    background-color: ${props => 
      props.isCorrect ? '#e6f7e6' : 
      props.isIncorrect ? '#ffebee' : 
      '#e3f2fd'
    };
  }
`;

const RadioInput = styled.input`
  margin-right: 0.75rem;
`;

const OptionText = styled.span`
  flex: 1;
`;

const Explanation = styled.div`
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 0.9rem;
  color: #666;
`;

const SubmitButton = styled.button`
  background-color: #16213e;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 1rem;
  
  &:hover {
    background-color: #0f3460;
  }
  
  &:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
`;

const ResultsContainer = styled.div`
  margin-top: 1rem;
  padding: 1rem;
  background-color: #e3f2fd;
  border-radius: 4px;
  text-align: center;
`;

const LoadingMessage = styled.div`
  padding: 1rem;
  text-align: center;
  color: #666;
`;

interface QuizQuestionsProps {
  sectionId: string;
}

const QuizQuestions: React.FC<QuizQuestionsProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { quizQuestions, isLoading } = useSelector((state: RootState) => state.books);
  const [selectedOptions, setSelectedOptions] = useState<Record<string, number>>({});
  const [submitted, setSubmitted] = useState(false);
  const [score, setScore] = useState<number | null>(null);
  
  useEffect(() => {
    dispatch(fetchQuizQuestions(sectionId));
  }, [dispatch, sectionId]);
  
  const handleOptionSelect = (questionId: string, optionIndex: number) => {
    if (submitted) return;
    
    setSelectedOptions({
      ...selectedOptions,
      [questionId]: optionIndex
    });
  };
  
  const handleSubmit = () => {
    if (Object.keys(selectedOptions).length === 0) return;
    
    let correctAnswers = 0;
    quizQuestions.forEach(question => {
      if (selectedOptions[question.id] === question.correctOptionIndex) {
        correctAnswers++;
      }
    });
    
    setScore(correctAnswers);
    setSubmitted(true);
  };
  
  const isOptionCorrect = (questionId: string, optionIndex: number) => {
    if (!submitted) return false;
    
    const question = quizQuestions.find(q => q.id === questionId);
    return question?.correctOptionIndex === optionIndex;
  };
  
  const isOptionIncorrect = (questionId: string, optionIndex: number) => {
    if (!submitted) return false;
    
    const question = quizQuestions.find(q => q.id === questionId);
    return selectedOptions[questionId] === optionIndex && question?.correctOptionIndex !== optionIndex;
  };
  
  if (isLoading) {
    return (
      <QuizContainer id="quiz-section">
        <Title>Quiz</Title>
        <LoadingMessage>Loading quiz questions...</LoadingMessage>
      </QuizContainer>
    );
  }
  
  if (quizQuestions.length === 0) {
    return (
      <QuizContainer id="quiz-section">
        <Title>Quiz</Title>
        <p>No quiz questions available for this section.</p>
      </QuizContainer>
    );
  }
  
  return (
    <QuizContainer id="quiz-section">
      <Title>Quiz</Title>
      
      <QuestionsList>
        {quizQuestions.map((question, questionIndex) => (
          <QuestionItem key={question.id}>
            <QuestionText>{questionIndex + 1}. {question.question}</QuestionText>
            
            <OptionsList>
              {question.options.map((option, optionIndex) => (
                <OptionItem 
                  key={optionIndex}
                  isSelected={selectedOptions[question.id] === optionIndex}
                  isCorrect={isOptionCorrect(question.id, optionIndex)}
                  isIncorrect={isOptionIncorrect(question.id, optionIndex)}
                >
                  <RadioInput 
                    type="radio"
                    name={`question-${question.id}`}
                    checked={selectedOptions[question.id] === optionIndex}
                    onChange={() => handleOptionSelect(question.id, optionIndex)}
                    disabled={submitted}
                  />
                  <OptionText>{option}</OptionText>
                </OptionItem>
              ))}
            </OptionsList>
            
            {submitted && question.explanation && (
              <Explanation>
                <strong>Explanation:</strong> {question.explanation}
              </Explanation>
            )}
          </QuestionItem>
        ))}
      </QuestionsList>
      
      {!submitted ? (
        <SubmitButton 
          onClick={handleSubmit}
          disabled={Object.keys(selectedOptions).length !== quizQuestions.length}
        >
          Submit Answers
        </SubmitButton>
      ) : (
        <ResultsContainer>
          <h3>Your Score: {score} out of {quizQuestions.length}</h3>
          <p>
            {score === quizQuestions.length 
              ? 'Perfect! You got all the answers correct.' 
              : `You got ${score} out of ${quizQuestions.length} questions correct.`}
          </p>
        </ResultsContainer>
      )}
    </QuizContainer>
  );
};

export default QuizQuestions;
