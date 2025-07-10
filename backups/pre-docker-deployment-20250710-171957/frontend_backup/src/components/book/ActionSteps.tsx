import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { fetchActionSteps } from '../../features/books/booksSlice';

const ActionStepsContainer = styled.div`
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

const StepsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const StepItem = styled.div`
  background-color: white;
  border-radius: 4px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
`;

const StepHeader = styled.div`
  display: flex;
  align-items: flex-start;
  margin-bottom: 0.5rem;
`;

const Checkbox = styled.input`
  margin-right: 0.75rem;
  margin-top: 0.25rem;
  width: 1.2rem;
  height: 1.2rem;
  cursor: pointer;
`;

const StepTitle = styled.h4`
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
  color: #16213e;
  flex: 1;
`;

const StepDescription = styled.p`
  color: #666;
  margin-bottom: 0.5rem;
  margin-left: 2rem;
`;

const StepPoints = styled.div`
  margin-left: 2rem;
  font-size: 0.9rem;
  color: #e94560;
  font-weight: bold;
`;

const LoadingMessage = styled.div`
  padding: 1rem;
  text-align: center;
  color: #666;
`;

interface ActionStepsProps {
  sectionId: string;
}

const ActionSteps: React.FC<ActionStepsProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { actionSteps, isLoading } = useSelector((state: RootState) => state.books);
  const [completedSteps, setCompletedSteps] = useState<Record<string, boolean>>({});
  
  useEffect(() => {
    dispatch(fetchActionSteps(sectionId));
  }, [dispatch, sectionId]);
  
  useEffect(() => {
    // Initialize completed steps from localStorage
    const savedSteps = localStorage.getItem(`completed_steps_${sectionId}`);
    if (savedSteps) {
      setCompletedSteps(JSON.parse(savedSteps));
    } else if (actionSteps.length > 0) {
      // Initialize with the completed status from the API
      const initialState: Record<string, boolean> = {};
      actionSteps.forEach(step => {
        initialState[step.id] = step.completed;
      });
      setCompletedSteps(initialState);
    }
  }, [actionSteps, sectionId]);
  
  const handleToggleStep = (stepId: string) => {
    const newCompletedSteps = {
      ...completedSteps,
      [stepId]: !completedSteps[stepId]
    };
    
    setCompletedSteps(newCompletedSteps);
    
    // Save to localStorage
    localStorage.setItem(`completed_steps_${sectionId}`, JSON.stringify(newCompletedSteps));
    
    // TODO: Update on the server if user is authenticated
  };
  
  if (isLoading) {
    return (
      <ActionStepsContainer id="action-steps">
        <Title>Action Steps</Title>
        <LoadingMessage>Loading action steps...</LoadingMessage>
      </ActionStepsContainer>
    );
  }
  
  if (actionSteps.length === 0) {
    return (
      <ActionStepsContainer id="action-steps">
        <Title>Action Steps</Title>
        <p>No action steps available for this section.</p>
      </ActionStepsContainer>
    );
  }
  
  return (
    <ActionStepsContainer id="action-steps">
      <Title>Action Steps</Title>
      
      <StepsList>
        {actionSteps.map((step) => (
          <StepItem key={step.id}>
            <StepHeader>
              <Checkbox 
                type="checkbox" 
                checked={completedSteps[step.id] || false}
                onChange={() => handleToggleStep(step.id)}
              />
              <StepTitle>{step.title}</StepTitle>
            </StepHeader>
            <StepDescription>{step.description}</StepDescription>
            <StepPoints>+{step.points} points when completed</StepPoints>
          </StepItem>
        ))}
      </StepsList>
    </ActionStepsContainer>
  );
};

export default ActionSteps;
