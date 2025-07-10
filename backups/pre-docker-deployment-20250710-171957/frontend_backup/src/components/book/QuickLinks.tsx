import React from 'react';
import styled from 'styled-components';

const QuickLinksContainer = styled.div`
  margin: 2rem 0;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border-bottom: 1px solid #ddd;
`;

const Title = styled.h4`
  font-size: 1.1rem;
  margin-bottom: 1rem;
  color: #16213e;
`;

const ButtonsContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
`;

const LinkButton = styled.button`
  background-color: #16213e;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-size: 0.9rem;
  
  &:hover {
    background-color: #0f3460;
  }
`;

interface QuickLinksProps {
  hasForumTopics?: boolean;
  hasActionSteps?: boolean;
  hasQuizQuestions?: boolean;
}

const QuickLinks: React.FC<QuickLinksProps> = ({
  hasForumTopics = false,
  hasActionSteps = false,
  hasQuizQuestions = false,
}) => {
  const scrollToSection = (sectionId: string) => {
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };
  
  return (
    <QuickLinksContainer>
      <Title>Quick Links</Title>
      <ButtonsContainer>
        {hasForumTopics && (
          <LinkButton onClick={() => scrollToSection('forum-topics')}>
            Forum Topics
          </LinkButton>
        )}
        
        {hasActionSteps && (
          <LinkButton onClick={() => scrollToSection('action-steps')}>
            Action Steps
          </LinkButton>
        )}
        
        {hasQuizQuestions && (
          <LinkButton onClick={() => scrollToSection('quiz-section')}>
            Quiz
          </LinkButton>
        )}
        
        <LinkButton onClick={() => scrollToSection('audio-book-section')}>
          Audio Book
        </LinkButton>
        
        <LinkButton onClick={() => scrollToSection('photo-book-section')}>
          Photo Book
        </LinkButton>
        
        <LinkButton onClick={() => scrollToSection('video-book-section')}>
          Video Book
        </LinkButton>
        
        <LinkButton onClick={() => scrollToSection('pdf-book-section')}>
          PDF Book
        </LinkButton>
      </ButtonsContainer>
    </QuickLinksContainer>
  );
};

export default QuickLinks;
