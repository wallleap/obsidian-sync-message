import React, { useState } from 'react';
import { useUserID } from '../hooks/useUserID';

interface AuthPageProps {
  onAuthenticated: (userID: string) => void;
}

export function AuthPage({ onAuthenticated }: AuthPageProps) {
  const { userID, loading, error, createUserID, setUserIDManually, downloadUserID } = useUserID();
  const [inputID, setInputID] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);
  const [showUploadOption, setShowUploadOption] = useState(false);

  const handleGenerate = async () => {
    setIsGenerating(true);
    try {
      const newID = await createUserID();
      if (newID) {
        onAuthenticated(newID);
      }
    } catch (err) {
      console.error('Failed to generate ID:', err);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleVerify = async () => {
    if (!inputID.trim()) return;
    
    setIsVerifying(true);
    try {
      const success = await setUserIDManually(inputID.trim());
      if (success) {
        onAuthenticated(inputID.trim());
      }
    } catch (err) {
      console.error('Failed to verify ID:', err);
    } finally {
      setIsVerifying(false);
    }
  };

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setInputID(content.trim());
    };
    reader.readAsText(file);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (userID) {
    onAuthenticated(userID);
    return null;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-indigo-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-800">OB Sync</h1>
          <p className="text-gray-500 mt-2">Connect your notes across devices</p>
        </div>

        {error && (
          <div className="mb-6 p-3 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
            {error}
          </div>
        )}

        <div className="space-y-4">
          <button
            onClick={handleGenerate}
            disabled={isGenerating}
            className="w-full py-3 px-4 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isGenerating ? (
              <span className="flex items-center justify-center">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Generating...
              </span>
            ) : (
              'Generate New User ID'
            )}
          </button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-white text-gray-500">Or</span>
            </div>
          </div>

          <div className="space-y-3">
            <input
              type="text"
              value={inputID}
              onChange={(e) => setInputID(e.target.value)}
              placeholder="Enter your User ID"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            />
            
            <div className="flex gap-2">
              <button
                onClick={() => setShowUploadOption(!showUploadOption)}
                className="flex-1 py-2 px-4 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-300 transition-colors text-sm"
              >
                Upload ID File
              </button>
              <button
                onClick={handleVerify}
                disabled={isVerifying || !inputID.trim()}
                className="flex-1 py-2 px-4 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm"
              >
                {isVerifying ? 'Verifying...' : 'Verify ID'}
              </button>
            </div>

            {showUploadOption && (
              <div className="mt-3 p-4 border border-dashed border-gray-300 rounded-lg text-center">
                <input
                  type="file"
                  accept=".txt"
                  onChange={handleFileUpload}
                  className="hidden"
                  id="file-upload"
                />
                <label
                  htmlFor="file-upload"
                  className="cursor-pointer block py-2 text-indigo-600 hover:text-indigo-700 font-medium"
                >
                  Click to upload your user ID file
                </label>
                <p className="text-xs text-gray-500 mt-1">Supports .txt files</p>
              </div>
            )}
          </div>
        </div>

        {userID && (
          <div className="mt-6 p-4 bg-indigo-50 rounded-lg">
            <p className="text-sm text-indigo-600 mb-2">Your User ID has been saved locally</p>
            <button
              onClick={downloadUserID}
              className="text-sm text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Download ID file
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
