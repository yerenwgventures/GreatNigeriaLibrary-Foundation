import React from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

const FooterContainer = styled.footer`
  background-color: #1a1a2e;
  color: white;
  padding: 3rem 0 1rem;
`;

const FooterContent = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
`;

const FooterSection = styled.div`
  margin-bottom: 2rem;
`;

const FooterTitle = styled.h3`
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: white;
`;

const FooterLinks = styled.ul`
  list-style: none;
  padding: 0;
  margin: 0;
`;

const FooterLink = styled.li`
  margin-bottom: 0.5rem;
`;

const StyledLink = styled(Link)`
  color: #e6e6e6;
  text-decoration: none;
  
  &:hover {
    color: white;
    text-decoration: underline;
  }
`;

const FooterBottom = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  margin-top: 2rem;
`;

const Copyright = styled.p`
  color: #e6e6e6;
  font-size: 0.9rem;
`;

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();
  
  return (
    <FooterContainer>
      <FooterContent>
        <FooterSection>
          <FooterTitle>Great Nigeria</FooterTitle>
          <p>Transforming Nigeria through citizen education, community building, and coordinated action.</p>
          <Copyright>© {currentYear} Great Nigeria. All rights reserved.</Copyright>
        </FooterSection>
        
        <FooterSection>
          <FooterTitle>Resources</FooterTitle>
          <FooterLinks>
            <FooterLink>
              <StyledLink to="/books">eBooks</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/resources">Implementation Tools</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/support-author">Support the Author</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/marketplace">Marketplace</StyledLink>
            </FooterLink>
          </FooterLinks>
        </FooterSection>
        
        <FooterSection>
          <FooterTitle>Quick Links</FooterTitle>
          <FooterLinks>
            <FooterLink>
              <StyledLink to="/about">About Us</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/community">Community</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/celebrate">Celebrate Nigeria</StyledLink>
            </FooterLink>
            <FooterLink>
              <StyledLink to="/contact">Contact Us</StyledLink>
            </FooterLink>
          </FooterLinks>
        </FooterSection>
      </FooterContent>
      
      <FooterBottom>
        <p>Built with ❤️ for Nigeria</p>
      </FooterBottom>
    </FooterContainer>
  );
};

export default Footer;
