import { useState, useEffect, useCallback } from 'react';
import localforage from 'localforage';
import { generateUserID, validateUserID } from '../api';

const STORAGE_KEY = 'ob_sync_user_id';

export function useUserID() {
  const [userID, setUserID] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadUserID = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const storedID = await localforage.getItem<string>(STORAGE_KEY);
      if (storedID) {
        const validation = await validateUserID(storedID);
        if (validation.valid) {
          setUserID(storedID);
        } else {
          setError('Invalid user ID, please generate a new one');
        }
      }
    } catch (err) {
      setError('Failed to load user ID');
      console.error('Error loading user ID:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const createUserID = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await generateUserID();
      await localforage.setItem(STORAGE_KEY, response.user_id);
      setUserID(response.user_id);
      return response.user_id;
    } catch (err) {
      setError('Failed to generate user ID');
      console.error('Error generating user ID:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const setUserIDManually = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const validation = await validateUserID(id);
      if (validation.valid) {
        await localforage.setItem(STORAGE_KEY, id);
        setUserID(id);
        return true;
      } else {
        setError('Invalid user ID');
        return false;
      }
    } catch (err) {
      setError('Failed to validate user ID');
      console.error('Error validating user ID:', err);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const downloadUserID = useCallback(() => {
    if (!userID) return;
    
    const blob = new Blob([userID], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'obsync_user_id.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [userID]);

  const clearUserID = useCallback(async () => {
    await localforage.removeItem(STORAGE_KEY);
    setUserID(null);
  }, []);

  useEffect(() => {
    loadUserID();
  }, [loadUserID]);

  return {
    userID,
    loading,
    error,
    createUserID,
    setUserIDManually,
    downloadUserID,
    clearUserID,
  };
}
