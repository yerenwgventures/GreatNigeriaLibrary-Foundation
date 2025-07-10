import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import styled from 'styled-components';
import { RootState } from '../store';
import { logout } from '../features/auth/authSlice';

const HeaderContainer = styled.header`
  background-color: #1a1a2e;
  color: white;
  padding: 1rem 0;
`;

const HeaderContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
`;

const Logo = styled(Link)`
  font-size: 1.5rem;
  font-weight: bold;
  color: white;
  text-decoration: none;

  &:hover {
    color: #e6e6e6;
  }
`;

const Nav = styled.nav`
  display: flex;
  align-items: center;
`;

const NavList = styled.ul`
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;

  @media (max-width: 768px) {
    display: none;
  }
`;

const NavItem = styled.li`
  margin-left: 1.5rem;
`;

const NavLink = styled(Link)<{ $active?: boolean }>`
  color: white;
  text-decoration: none;
  font-weight: ${(props) => (props.$active ? 'bold' : 'normal')};
  padding-bottom: 0.25rem;
  border-bottom: ${(props) => (props.$active ? '2px solid white' : 'none')};

  &:hover {
    color: #e6e6e6;
  }
`;

const AuthButtons = styled.div`
  display: flex;
  align-items: center;
  margin-left: 1.5rem;
`;

const Button = styled.button`
  background-color: transparent;
  color: white;
  border: 1px solid white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  margin-left: 0.5rem;

  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
`;

const PrimaryButton = styled(Button)`
  background-color: #16213e;
  border: none;

  &:hover {
    background-color: #0f3460;
  }
`;

const MobileMenuButton = styled.button`
  display: none;
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;

  @media (max-width: 768px) {
    display: block;
  }
`;

const MobileMenu = styled.div<{ isOpen: boolean }>`
  display: ${(props) => (props.isOpen ? 'block' : 'none')};
  position: absolute;
  top: 70px;
  left: 0;
  right: 0;
  background-color: #1a1a2e;
  padding: 1rem;
  z-index: 100;

  @media (min-width: 769px) {
    display: none;
  }
`;

const MobileNavList = styled.ul`
  list-style: none;
  margin: 0;
  padding: 0;
`;

const MobileNavItem = styled.li`
  margin: 1rem 0;
`;

const MobileNavLink = styled(Link)<{ $active?: boolean }>`
  color: white;
  text-decoration: none;
  font-weight: ${(props) => (props.$active ? 'bold' : 'normal')};
  display: block;
  padding: 0.5rem 0;

  &:hover {
    color: #e6e6e6;
  }
`;

const Header: React.FC = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(logout());
    navigate('/');
  };

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  const closeMobileMenu = () => {
    setMobileMenuOpen(false);
  };

  return (
    <HeaderContainer>
      <HeaderContent>
        <Logo to="/">Great Nigeria</Logo>

        <MobileMenuButton onClick={toggleMobileMenu}>
          {mobileMenuOpen ? '✕' : '☰'}
        </MobileMenuButton>

        <Nav>
          <NavList>
            <NavItem>
              <NavLink to="/">Home</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/books">eBooks</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/community">Community</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/groups">Local Groups</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/celebrate">Celebrate Nigeria</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/about">About</NavLink>
            </NavItem>
            <NavItem>
              <NavLink to="/resources">Resources</NavLink>
            </NavItem>
          </NavList>

          <AuthButtons>
            {isAuthenticated ? (
              <>
                <NavLink to="/profile">
                  {user?.name || 'Profile'}
                </NavLink>
                <Button onClick={handleLogout}>Logout</Button>
              </>
            ) : (
              <>
                <Button onClick={() => navigate('/login')}>Login</Button>
                <PrimaryButton onClick={() => navigate('/register')}>Register</PrimaryButton>
              </>
            )}
          </AuthButtons>
        </Nav>
      </HeaderContent>

      <MobileMenu isOpen={mobileMenuOpen}>
        <MobileNavList>
          <MobileNavItem>
            <MobileNavLink to="/" onClick={closeMobileMenu}>Home</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/books" onClick={closeMobileMenu}>eBooks</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/community" onClick={closeMobileMenu}>Community</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/groups" onClick={closeMobileMenu}>Local Groups</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/celebrate" onClick={closeMobileMenu}>Celebrate Nigeria</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/about" onClick={closeMobileMenu}>About</MobileNavLink>
          </MobileNavItem>
          <MobileNavItem>
            <MobileNavLink to="/resources" onClick={closeMobileMenu}>Resources</MobileNavLink>
          </MobileNavItem>

          {isAuthenticated ? (
            <>
              <MobileNavItem>
                <MobileNavLink to="/profile" onClick={closeMobileMenu}>Profile</MobileNavLink>
              </MobileNavItem>
              <MobileNavItem>
                <Button onClick={() => { handleLogout(); closeMobileMenu(); }}>Logout</Button>
              </MobileNavItem>
            </>
          ) : (
            <>
              <MobileNavItem>
                <Button onClick={() => { navigate('/login'); closeMobileMenu(); }}>Login</Button>
              </MobileNavItem>
              <MobileNavItem>
                <PrimaryButton onClick={() => { navigate('/register'); closeMobileMenu(); }}>Register</PrimaryButton>
              </MobileNavItem>
            </>
          )}
        </MobileNavList>
      </MobileMenu>
    </HeaderContainer>
  );
};

export default Header;
