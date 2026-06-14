import { useState } from 'react';
import { AuthPage } from './components/AuthPage';
import { MainPage } from './components/MainPage';

function App() {
  const [authenticatedUserID, setAuthenticatedUserID] = useState<string | null>(null);

  const handleAuthenticated = (userID: string) => {
    setAuthenticatedUserID(userID);
  };

  const handleLogout = () => {
    setAuthenticatedUserID(null);
  };

  if (!authenticatedUserID) {
    return <AuthPage onAuthenticated={handleAuthenticated} />;
  }

  return <MainPage userID={authenticatedUserID} onLogout={handleLogout} />;
}

export default App;
