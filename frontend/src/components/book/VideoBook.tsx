import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { generateVideoBook, getShareableLink } from '../../features/books/booksSlice';

const VideoBookContainer = styled.div`
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

const VideoPlayer = styled.video`
  width: 100%;
  max-height: 400px;
  margin: 1rem 0;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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

interface VideoBookProps {
  sectionId: string;
}

const VideoBook: React.FC<VideoBookProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { videoBook, isLoadingVideo } = useSelector((state: RootState) => state.books);
  const [isSharing, setIsSharing] = useState(false);
  
  const handleGenerateVideo = () => {
    dispatch(generateVideoBook(sectionId));
  };
  
  const handleShare = async () => {
    if (!videoBook) return;
    
    setIsSharing(true);
    try {
      await dispatch(getShareableLink({ sectionId, mediaType: 'video' }));
      
      // Use Web Share API if available
      if (navigator.share) {
        await navigator.share({
          title: videoBook.title,
          text: `Watch this video: ${videoBook.title}`,
          url: videoBook.videoUrl,
        });
      } else {
        // Fallback to copying to clipboard
        await navigator.clipboard.writeText(videoBook.videoUrl);
        alert('Link copied to clipboard!');
      }
    } catch (error) {
      console.error('Error sharing:', error);
    } finally {
      setIsSharing(false);
    }
  };
  
  return (
    <VideoBookContainer id="video-book-section">
      <Title>Video Book</Title>
      
      {videoBook ? (
        <>
          <VideoPlayer 
            controls 
            src={videoBook.videoUrl} 
            poster={videoBook.thumbnailUrl}
          >
            Your browser does not support the video element.
          </VideoPlayer>
          
          <div>
            <p>Duration: {Math.floor(videoBook.duration / 60)}:{(videoBook.duration % 60).toString().padStart(2, '0')}</p>
            <p>Generated: {new Date(videoBook.generatedAt).toLocaleString()}</p>
          </div>
          
          <ShareButton onClick={handleShare} disabled={isSharing}>
            {isSharing && <LoadingSpinner />}
            {isSharing ? 'Sharing...' : 'Share Video'}
          </ShareButton>
        </>
      ) : (
        <Button onClick={handleGenerateVideo} disabled={isLoadingVideo}>
          {isLoadingVideo && <LoadingSpinner />}
          {isLoadingVideo ? 'Generating Video...' : 'Generate Video Book'}
        </Button>
      )}
    </VideoBookContainer>
  );
};

export default VideoBook;
