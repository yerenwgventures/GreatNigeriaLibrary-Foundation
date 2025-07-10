import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../../store';
import { generatePdfBook, getShareableLink } from '../../features/books/booksSlice';

const PDFBookContainer = styled.div`
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

const PDFPreview = styled.div`
  margin: 1rem 0;
  padding: 1rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const PDFIcon = styled.div`
  width: 40px;
  height: 50px;
  background-color: #e94560;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  margin-right: 1rem;
`;

const PDFInfo = styled.div`
  flex: 1;
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

const DownloadButton = styled(Button)`
  background-color: #4caf50;
  
  &:hover {
    background-color: #3e8e41;
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

interface PDFBookProps {
  sectionId: string;
}

const PDFBook: React.FC<PDFBookProps> = ({ sectionId }) => {
  const dispatch = useDispatch();
  const { pdfBook, isLoadingPdf } = useSelector((state: RootState) => state.books);
  const [isSharing, setIsSharing] = useState(false);
  
  const handleGeneratePdf = () => {
    dispatch(generatePdfBook(sectionId));
  };
  
  const handleShare = async () => {
    if (!pdfBook) return;
    
    setIsSharing(true);
    try {
      await dispatch(getShareableLink({ sectionId, mediaType: 'pdf' }));
      
      // Use Web Share API if available
      if (navigator.share) {
        await navigator.share({
          title: pdfBook.title,
          text: `Check out this PDF: ${pdfBook.title}`,
          url: pdfBook.pdfUrl,
        });
      } else {
        // Fallback to copying to clipboard
        await navigator.clipboard.writeText(pdfBook.pdfUrl);
        alert('Link copied to clipboard!');
      }
    } catch (error) {
      console.error('Error sharing:', error);
    } finally {
      setIsSharing(false);
    }
  };
  
  return (
    <PDFBookContainer id="pdf-book-section">
      <Title>PDF Book</Title>
      
      {pdfBook ? (
        <>
          <PDFPreview>
            <PDFIcon>PDF</PDFIcon>
            <PDFInfo>
              <h4>{pdfBook.title}</h4>
              <p>{pdfBook.pageCount} pages • Generated: {new Date(pdfBook.generatedAt).toLocaleString()}</p>
            </PDFInfo>
            <DownloadButton 
              as="a" 
              href={pdfBook.pdfUrl} 
              target="_blank" 
              rel="noopener noreferrer"
            >
              Download
            </DownloadButton>
          </PDFPreview>
          
          <ShareButton onClick={handleShare} disabled={isSharing}>
            {isSharing && <LoadingSpinner />}
            {isSharing ? 'Sharing...' : 'Share PDF'}
          </ShareButton>
        </>
      ) : (
        <Button onClick={handleGeneratePdf} disabled={isLoadingPdf}>
          {isLoadingPdf && <LoadingSpinner />}
          {isLoadingPdf ? 'Generating PDF...' : 'Generate PDF Book'}
        </Button>
      )}
    </PDFBookContainer>
  );
};

export default PDFBook;
