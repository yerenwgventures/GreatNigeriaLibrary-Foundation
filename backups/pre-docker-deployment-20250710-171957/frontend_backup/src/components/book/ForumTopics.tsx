import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { Link } from 'react-router-dom';
import { RootState } from '../../store';
import { fetchForumTopics } from '../../features/books/booksSlice';

const ForumTopicsContainer = styled.div`
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

const TopicsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const TopicItem = styled.div`
  background-color: white;
  border-radius: 4px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
`;

const TopicTitle = styled.h4`
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
  color: #16213e;
`;

const TopicDescription = styled.p`
  color: #666;
  margin-bottom: 1rem;
`;

const TopicFooter = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
  color: #888;
`;

const ViewButton = styled(Link)`
  background-color: #16213e;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.9rem;
  
  &:hover {
    background-color: #0f3460;
  }
`;

const LoadingMessage = styled.div`
  padding: 1rem;
  text-align: center;
  color: #666;
`;

interface ForumTopicsProps {
  sectionId: string;
}

const ForumTopics: React.FC<ForumTopicsProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { forumTopics, isLoading } = useSelector((state: RootState) => state.books);
  
  useEffect(() => {
    dispatch(fetchForumTopics(sectionId));
  }, [dispatch, sectionId]);
  
  if (isLoading) {
    return (
      <ForumTopicsContainer id="forum-topics">
        <Title>Forum Topics</Title>
        <LoadingMessage>Loading forum topics...</LoadingMessage>
      </ForumTopicsContainer>
    );
  }
  
  if (forumTopics.length === 0) {
    return (
      <ForumTopicsContainer id="forum-topics">
        <Title>Forum Topics</Title>
        <p>No forum topics available for this section.</p>
      </ForumTopicsContainer>
    );
  }
  
  return (
    <ForumTopicsContainer id="forum-topics">
      <Title>Forum Topics</Title>
      
      <TopicsList>
        {forumTopics.map((topic) => (
          <TopicItem key={topic.id}>
            <TopicTitle>{topic.title}</TopicTitle>
            <TopicDescription>{topic.description}</TopicDescription>
            <TopicFooter>
              <span>{topic.responseCount} responses</span>
              <ViewButton to={`/forum/topics/${topic.id}`}>
                View Discussion
              </ViewButton>
            </TopicFooter>
          </TopicItem>
        ))}
      </TopicsList>
    </ForumTopicsContainer>
  );
};

export default ForumTopics;
