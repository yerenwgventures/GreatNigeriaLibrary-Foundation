import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  FormControlLabel,
  Grid,
  Paper,
  Radio,
  RadioGroup,
  Step,
  StepLabel,
  Stepper,
  Typography,
  useTheme,
} from '@mui/material';
import {
  fetchAssessmentQuestions,
  submitAssessment,
  selectAssessmentQuestions,
  selectAssessmentResult,
  selectPersonalizationLoading,
  selectPersonalizationError,
  clearAssessmentResult,
} from '../../features/personalization/personalizationSlice';
import { AssessmentQuestion, AssessmentResponse } from '../../api/personalizationService';
import { AppDispatch } from '../../store';
import { useNavigate } from 'react-router-dom';
import { styled } from '@mui/system';

// Styled components
const AssessmentCard = styled(Card)(({ theme }) => ({
  marginBottom: theme.spacing(3),
  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
  borderRadius: theme.spacing(2),
  overflow: 'hidden',
}));

const QuestionContainer = styled(Box)(({ theme }) => ({
  padding: theme.spacing(3),
}));

const OptionContainer = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
}));

const ResultCard = styled(Card)(({ theme }) => ({
  marginBottom: theme.spacing(3),
  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
  borderRadius: theme.spacing.unit,
  overflow: 'hidden',
}));

const StyleBar = styled(Box)<{ value: number; color: string }>(({ theme, value, color }) => ({
  height: 20,
  width: `${value}%`,
  backgroundColor: color,
  borderRadius: 10,
  transition: 'width 1s ease-in-out',
}));

// Main component
const LearningStyleAssessment: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const theme = useTheme();
  
  const questions = useSelector(selectAssessmentQuestions);
  const result = useSelector(selectAssessmentResult);
  const loading = useSelector(selectPersonalizationLoading);
  const error = useSelector(selectPersonalizationError);
  
  const [activeStep, setActiveStep] = useState(0);
  const [responses, setResponses] = useState<Record<number, number>>({});
  const [currentPage, setCurrentPage] = useState(0);
  const [showResults, setShowResults] = useState(false);
  const [resultDialogOpen, setResultDialogOpen] = useState(false);
  
  const questionsPerPage = 5;
  const totalPages = Math.ceil(questions.length / questionsPerPage);
  const currentQuestions = questions.slice(
    currentPage * questionsPerPage,
    (currentPage + 1) * questionsPerPage
  );
  
  useEffect(() => {
    dispatch(fetchAssessmentQuestions());
    
    // Clean up assessment result when component unmounts
    return () => {
      dispatch(clearAssessmentResult());
    };
  }, [dispatch]);
  
  useEffect(() => {
    if (result) {
      setShowResults(true);
      setResultDialogOpen(true);
    }
  }, [result]);
  
  const handleOptionSelect = (questionId: number, optionIndex: number) => {
    setResponses({
      ...responses,
      [questionId]: optionIndex,
    });
  };
  
  const handleNext = () => {
    if (currentPage < totalPages - 1) {
      setCurrentPage(currentPage + 1);
      setActiveStep(activeStep + 1);
    } else {
      handleSubmit();
    }
  };
  
  const handleBack = () => {
    if (currentPage > 0) {
      setCurrentPage(currentPage - 1);
      setActiveStep(activeStep - 1);
    }
  };
  
  const handleSubmit = () => {
    const assessmentResponses: AssessmentResponse[] = Object.keys(responses).map((questionId) => ({
      questionId: parseInt(questionId),
      selectedOption: responses[parseInt(questionId)],
    }));
    
    dispatch(submitAssessment(assessmentResponses));
  };
  
  const handleCloseResultDialog = () => {
    setResultDialogOpen(false);
  };
  
  const handleViewPaths = () => {
    setResultDialogOpen(false);
    navigate('/personalized-paths');
  };
  
  const isNextDisabled = () => {
    return currentQuestions.some((question) => responses[question.id] === undefined);
  };
  
  const getStyleColor = (style: string) => {
    const colors: Record<string, string> = {
      visual: theme.palette.primary.main,
      auditory: theme.palette.secondary.main,
      readWrite: theme.palette.success.main,
      kinesthetic: theme.palette.warning.main,
      social: theme.palette.info.main,
      solitary: theme.palette.error.main,
      logical: theme.palette.grey[700],
    };
    
    return colors[style] || theme.palette.primary.main;
  };
  
  if (loading && questions.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }
  
  if (error) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <Typography color="error" variant="h6">
          Error: {error}
        </Typography>
      </Box>
    );
  }
  
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Paper sx={{ p: 3, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          Learning Style Assessment
        </Typography>
        <Typography variant="body1" paragraph align="center">
          Discover your unique learning style by answering the following questions. This will help us personalize your
          learning experience.
        </Typography>
        
        <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 4 }}>
          {Array.from({ length: totalPages }).map((_, index) => (
            <Step key={index}>
              <StepLabel>Page {index + 1}</StepLabel>
            </Step>
          ))}
        </Stepper>
        
        {currentQuestions.map((question: AssessmentQuestion) => (
          <AssessmentCard key={question.id}>
            <QuestionContainer>
              <Typography variant="h6" gutterBottom>
                {question.question}
              </Typography>
              <FormControl component="fieldset" fullWidth>
                <RadioGroup
                  value={responses[question.id] !== undefined ? responses[question.id] : ''}
                  onChange={(e) => handleOptionSelect(question.id, parseInt(e.target.value))}
                >
                  {question.options.map((option, index) => (
                    <FormControlLabel
                      key={index}
                      value={index}
                      control={<Radio />}
                      label={option}
                      sx={{ py: 1 }}
                    />
                  ))}
                </RadioGroup>
              </FormControl>
            </QuestionContainer>
          </AssessmentCard>
        ))}
        
        <Box display="flex" justifyContent="space-between" mt={3}>
          <Button
            variant="outlined"
            onClick={handleBack}
            disabled={currentPage === 0}
          >
            Back
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={handleNext}
            disabled={isNextDisabled()}
          >
            {currentPage < totalPages - 1 ? 'Next' : 'Submit'}
          </Button>
        </Box>
      </Paper>
      
      {/* Results Dialog */}
      <Dialog
        open={resultDialogOpen}
        onClose={handleCloseResultDialog}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Your Learning Style Results</DialogTitle>
        <DialogContent>
          {result && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Your Primary Learning Style: <strong>{result.learningStyle.primaryStyle}</strong>
              </Typography>
              <Typography variant="body1" paragraph>
                Secondary Style: {result.learningStyle.secondaryStyle}
              </Typography>
              
              <Grid container spacing={2} sx={{ mb: 3 }}>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Visual: {result.learningStyle.visual}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.visual} color={getStyleColor('visual')} />
                  </Box>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Auditory: {result.learningStyle.auditory}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.auditory} color={getStyleColor('auditory')} />
                  </Box>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Read/Write: {result.learningStyle.readWrite}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.readWrite} color={getStyleColor('readWrite')} />
                  </Box>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Kinesthetic: {result.learningStyle.kinesthetic}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.kinesthetic} color={getStyleColor('kinesthetic')} />
                  </Box>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Social: {result.learningStyle.social}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.social} color={getStyleColor('social')} />
                  </Box>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Typography variant="subtitle2">Solitary: {result.learningStyle.solitary}%</Typography>
                  <Box sx={{ bgcolor: 'grey.300', borderRadius: 10, mb: 1 }}>
                    <StyleBar value={result.learningStyle.solitary} color={getStyleColor('solitary')} />
                  </Box>
                </Grid>
              </Grid>
              
              <Divider sx={{ my: 2 }} />
              
              <Typography variant="h6" gutterBottom>
                What This Means For You
              </Typography>
              
              {result.learningStyle.primaryStyle === 'visual' && (
                <Typography variant="body1" paragraph>
                  As a visual learner, you learn best through seeing. You benefit from diagrams, charts, videos, and written instructions. Try using color-coding, mind maps, and flashcards in your studies.
                </Typography>
              )}
              
              {result.learningStyle.primaryStyle === 'auditory' && (
                <Typography variant="body1" paragraph>
                  As an auditory learner, you learn best through listening. You benefit from lectures, discussions, and audio materials. Try reading aloud, participating in group discussions, and using recorded materials.
                </Typography>
              )}
              
              {result.learningStyle.primaryStyle === 'readWrite' && (
                <Typography variant="body1" paragraph>
                  As a read/write learner, you learn best through text. You benefit from reading and writing. Try taking detailed notes, rewriting information in your own words, and using written summaries.
                </Typography>
              )}
              
              {result.learningStyle.primaryStyle === 'kinesthetic' && (
                <Typography variant="body1" paragraph>
                  As a kinesthetic learner, you learn best through doing. You benefit from hands-on activities and practical experiences. Try role-playing, building models, and incorporating movement into your learning.
                </Typography>
              )}
              
              {result.learningStyle.primaryStyle === 'social' && (
                <Typography variant="body1" paragraph>
                  As a social learner, you learn best in groups. You benefit from collaboration and discussion. Try study groups, peer teaching, and collaborative projects.
                </Typography>
              )}
              
              {result.learningStyle.primaryStyle === 'solitary' && (
                <Typography variant="body1" paragraph>
                  As a solitary learner, you learn best alone. You benefit from self-study and reflection. Try independent research, journaling, and setting personal goals.
                </Typography>
              )}
              
              <Typography variant="body1" paragraph>
                Based on your learning style, we've created personalized learning paths and recommendations for you. Click "View Personalized Paths" to explore them.
              </Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseResultDialog}>Close</Button>
          <Button onClick={handleViewPaths} variant="contained" color="primary">
            View Personalized Paths
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default LearningStyleAssessment;
