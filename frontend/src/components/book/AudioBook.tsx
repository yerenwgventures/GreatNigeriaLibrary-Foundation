import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { generateAudioBook, getShareableLink } from '../../features/books/booksSlice';

const AudioBookContainer = styled.div`
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

const AudioPlayer = styled.audio`
  width: 100%;
  margin: 1rem 0;
`;

const Button = styled.button`
  background-color: #16213e;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover {
    background-color: #0f3460;
  }
  
  &:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
`;

const ShareButton = styled(Button)`
  background-color: #e94560;
  margin-top: 1rem;
  
  &:hover {
    background-color: #d13050;
  }
`;

const LoadingSpinner = styled.div`
  display: inline-block;
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
  
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
`;

interface AudioBookProps {
  sectionId: string;
}

const AudioBook: React.FC<AudioBookProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { audioBook, isLoadingAudio } = useSelector((state: RootState) => state.books);
  const [isSharing, setIsSharing] = useState(false);
  
  const handleGenerateAudio = () => {
    dispatch(generateAudioBook(sectionId));
  };
  
  const handleShare = async () => {
    if (!audioBook) return;
    
    setIsSharing(true);
    try {
      await dispatch(getShareableLink({ sectionId, mediaType: 'audio' }));
      
      // Use Web Share API if available
      if (navigator.share) {
        await navigator.share({
          title: audioBook.title,
          text: `Listen to this audio book: ${audioBook.title}`,
          url: audioBook.audioUrl,
        });
      } else {
        // Fallback to copying to clipboard
        await navigator.clipboard.writeText(audioBook.audioUrl);
        alert('Link copied to clipboard!');
      }
    } catch (error) {
      console.error('Error sharing:', error);
    } finally {
      setIsSharing(false);
    }
  };
  
  return (
    <AudioBookContainer id="audio-book-section">
      <Title>Audio Book</Title>
      
      {audioBook ? (
        <>
          <AudioPlayer controls src={audioBook.audioUrl}>
            Your browser does not support the audio element.
          </AudioPlayer>
          
          <div>
            <p>Duration: {Math.floor(audioBook.duration / 60)}:{(audioBook.duration % 60).toString().padStart(2, '0')}</p>
            <p>Generated: {new Date(audioBook.generatedAt).toLocaleString()}</p>
          </div>
          
          <ShareButton onClick={handleShare} disabled={isSharing}>
            {isSharing && <LoadingSpinner />}
            {isSharing ? 'Sharing...' : 'Share Audio'}
          </ShareButton>
        </>
      ) : (
        <Button onClick={handleGenerateAudio} disabled={isLoadingAudio}>
          {isLoadingAudio && <LoadingSpinner />}
          {isLoadingAudio ? 'Generating Audio...' : 'Generate Audio Book'}
        </Button>
      )}
    </AudioBookContainer>
  );
};

export default AudioBook;
