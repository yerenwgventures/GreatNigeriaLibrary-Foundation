import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { generatePhotoBook, getShareableLink } from '../../features/books/booksSlice';

const PhotoBookContainer = styled.div`
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

const PhotoGallery = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  margin: 1rem 0;
`;

const Photo = styled.div`
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  
  img {
    width: 100%;
    height: 200px;
    object-fit: cover;
    transition: transform 0.3s ease;
    
    &:hover {
      transform: scale(1.05);
    }
  }
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

interface PhotoBookProps {
  sectionId: string;
}

const PhotoBook: React.FC<PhotoBookProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { photoBook, isLoadingPhoto } = useSelector((state: RootState) => state.books);
  const [isSharing, setIsSharing] = useState(false);
  
  const handleGeneratePhotoBook = () => {
    dispatch(generatePhotoBook(sectionId));
  };
  
  const handleShare = async () => {
    if (!photoBook) return;
    
    setIsSharing(true);
    try {
      await dispatch(getShareableLink({ sectionId, mediaType: 'photo' }));
      
      // Use Web Share API if available
      if (navigator.share) {
        await navigator.share({
          title: photoBook.title,
          text: `Check out this photo collection: ${photoBook.title}`,
          url: window.location.href,
        });
      } else {
        // Fallback to copying to clipboard
        await navigator.clipboard.writeText(window.location.href);
        alert('Link copied to clipboard!');
      }
    } catch (error) {
      console.error('Error sharing:', error);
    } finally {
      setIsSharing(false);
    }
  };
  
  return (
    <PhotoBookContainer id="photo-book-section">
      <Title>Photo Book</Title>
      
      {photoBook ? (
        <>
          <PhotoGallery>
            {photoBook.photoUrls.map((url, index) => (
              <Photo key={index}>
                <img src={url} alt={`Illustration ${index + 1}`} />
              </Photo>
            ))}
          </PhotoGallery>
          
          <div>
            <p>Photos: {photoBook.count}</p>
            <p>Generated: {new Date(photoBook.generatedAt).toLocaleString()}</p>
          </div>
          
          <ShareButton onClick={handleShare} disabled={isSharing}>
            {isSharing && <LoadingSpinner />}
            {isSharing ? 'Sharing...' : 'Share Photo Collection'}
          </ShareButton>
        </>
      ) : (
        <Button onClick={handleGeneratePhotoBook} disabled={isLoadingPhoto}>
          {isLoadingPhoto && <LoadingSpinner />}
          {isLoadingPhoto ? 'Generating Photos...' : 'Generate Photo Book'}
        </Button>
      )}
    </PhotoBookContainer>
  );
};

export default PhotoBook;
