import React, { useState, useCallback, useEffect } from 'react';
import { Message, sendMessage, uploadAttachment, syncMessages } from '../api';
import { useUserID } from '../hooks/useUserID';

interface MainPageProps {
  userID: string;
  onLogout: () => void;
}

type ProcessingStage = 'idle' | 'sending' | 'fetching' | 'converting' | 'syncing' | 'success' | 'error';

export function MainPage({ userID, onLogout }: MainPageProps) {
  const { downloadUserID, clearUserID } = useUserID();
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputText, setInputText] = useState('');
  const [inputURL, setInputURL] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  const [sendStatus, setSendStatus] = useState<'success' | 'error' | null>(null);
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [processingStage, setProcessingStage] = useState<ProcessingStage>('idle');
  const [showProcessingOverlay, setShowProcessingOverlay] = useState(false);
  const [manualURL, setManualURL] = useState<string | null>(null);

  const handleSend = async () => {
    if (!inputText.trim() && !inputURL.trim() && !uploadFile) return;
    
    setIsSending(true);
    setSendStatus(null);
    setShowProcessingOverlay(true);
    
    try {
      // Stage 1: Sending
      setProcessingStage('sending');
      await new Promise(resolve => setTimeout(resolve, 300)); // Brief pause for UI
      
      if (uploadFile) {
        await uploadAttachment(userID, uploadFile);
      } else if (inputURL.trim()) {
        setProcessingStage('fetching');
        const response = await sendMessage(userID, 'url', inputURL.trim(), inputURL.trim());
        
        // Check if manual processing is needed
        if (response.needs_manual) {
          setManualURL(response.original_url || inputURL.trim());
          setShowProcessingOverlay(false);
          setProcessingStage('idle');
          setIsSending(false);
          return;
        }
      } else {
        await sendMessage(userID, 'text', inputText.trim());
      }
      
      // Stage 2: Syncing
      setProcessingStage('syncing');
      await handleSync();
      
      setProcessingStage('success');
      setSendStatus('success');
      setInputText('');
      setInputURL('');
      setUploadFile(null);
      
      setTimeout(() => {
        setSendStatus(null);
        setShowProcessingOverlay(false);
        setProcessingStage('idle');
      }, 2000);
    } catch (err) {
      console.error('Failed to send:', err);
      setProcessingStage('error');
      setSendStatus('error');
      setTimeout(() => {
        setSendStatus(null);
        setShowProcessingOverlay(false);
        setProcessingStage('idle');
      }, 3000);
    } finally {
      setIsSending(false);
    }
  };

  const handleManualPaste = async (content: string) => {
    if (!content.trim()) return;
    
    setShowProcessingOverlay(true);
    setProcessingStage('sending');
    
    try {
      await sendMessage(userID, 'text', content.trim());
      setManualURL(null);
      setInputText('');
      setInputURL('');
      
      // Sync to see the new message
      setProcessingStage('syncing');
      await handleSync();
      
      setProcessingStage('success');
      setSendStatus('success');
      
      setTimeout(() => {
        setSendStatus(null);
        setShowProcessingOverlay(false);
        setProcessingStage('idle');
      }, 2000);
    } catch (err) {
      console.error('Failed to send manual content:', err);
      setProcessingStage('error');
      setSendStatus('error');
      setTimeout(() => {
        setSendStatus(null);
        setShowProcessingOverlay(false);
        setProcessingStage('idle');
      }, 3000);
    }
  };

  const handleSync = useCallback(async () => {
    setIsSyncing(true);
    try {
      const lastSync = localStorage.getItem('last_sync_time') || '';
      const newMessages = await syncMessages(userID, lastSync);
      
      if (newMessages.length > 0) {
        setMessages((prev) => {
          const existingIds = new Set(prev.map((m) => m.id));
          const filtered = newMessages.filter((m) => !existingIds.has(m.id));
          return [...prev, ...filtered];
        });
        localStorage.setItem('last_sync_time', new Date().toISOString());
      }
    } catch (err) {
      console.error('Failed to sync:', err);
    } finally {
      setIsSyncing(false);
    }
  }, [userID]);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setUploadFile(file);
    }
  };

  const handleLogout = async () => {
    await clearUserID();
    onLogout();
  };

  useEffect(() => {
    handleSync();
  }, [handleSync]);

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'text': return 'Text';
      case 'url': return 'URL';
      case 'attachment': return 'Attachment';
      default: return type;
    }
  };

  const getProcessingMessage = () => {
    switch (processingStage) {
      case 'sending': return 'Sending message...';
      case 'fetching': return 'Fetching URL content...';
      case 'converting': return 'Converting to markdown...';
      case 'syncing': return 'Syncing with server...';
      case 'success': return 'Completed!';
      case 'error': return 'Failed';
      default: return '';
    }
  };

  const getProcessingProgress = () => {
    switch (processingStage) {
      case 'sending': return 25;
      case 'fetching': return 40;
      case 'converting': return 60;
      case 'syncing': return 80;
      case 'success': return 100;
      case 'error': return 100;
      default: return 0;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Processing Overlay */}
      {showProcessingOverlay && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-2xl p-8 w-96 max-w-[90vw]">
            <div className="text-center">
              <div className="mb-6">
                {processingStage === 'success' ? (
                  <div className="w-16 h-16 mx-auto bg-green-100 rounded-full flex items-center justify-center">
                    <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                ) : processingStage === 'error' ? (
                  <div className="w-16 h-16 mx-auto bg-red-100 rounded-full flex items-center justify-center">
                    <svg className="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </div>
                ) : (
                  <div className="w-16 h-16 mx-auto bg-indigo-100 rounded-full flex items-center justify-center">
                    <svg className="w-8 h-8 text-indigo-600 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>
                )}
              </div>
              
              <h3 className="text-lg font-semibold text-gray-800 mb-2">
                {processingStage === 'success' ? 'Success' : processingStage === 'error' ? 'Error' : 'Processing'}
              </h3>
              
              <p className="text-gray-500 mb-6">{getProcessingMessage()}</p>
              
              {/* Progress Bar */}
              {processingStage !== 'idle' && (
                <div className="w-full bg-gray-200 rounded-full h-2 mb-4">
                  <div 
                    className={`h-2 rounded-full transition-all duration-500 ${
                      processingStage === 'error' ? 'bg-red-500' : 
                      processingStage === 'success' ? 'bg-green-500' : 'bg-indigo-500'
                    }`}
                    style={{ width: `${getProcessingProgress()}%` }}
                  ></div>
                </div>
              )}
              
              {/* Processing Steps */}
              <div className="space-y-2 text-sm">
                <div className={`flex items-center gap-2 ${
                  ['sending', 'fetching', 'converting', 'syncing', 'success'].includes(processingStage) 
                    ? 'text-indigo-600' : 'text-gray-400'
                }`}>
                  <div className={`w-5 h-5 rounded-full flex items-center justify-center ${
                    ['sending', 'fetching', 'converting', 'syncing', 'success'].includes(processingStage)
                      ? 'bg-indigo-100' : 'bg-gray-100'
                  }`}>
                    {['fetching', 'converting', 'syncing', 'success'].includes(processingStage) ? (
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    ) : processingStage === 'sending' ? (
                      <div className="w-2 h-2 bg-indigo-600 rounded-full animate-pulse"></div>
                    ) : (
                      <div className="w-2 h-2 bg-gray-300 rounded-full"></div>
                    )}
                  </div>
                  <span>Sending message</span>
                </div>
                
                <div className={`flex items-center gap-2 ${
                  ['fetching', 'converting', 'syncing', 'success'].includes(processingStage) 
                    ? 'text-indigo-600' : 'text-gray-400'
                }`}>
                  <div className={`w-5 h-5 rounded-full flex items-center justify-center ${
                    ['fetching', 'converting', 'syncing', 'success'].includes(processingStage)
                      ? 'bg-indigo-100' : 'bg-gray-100'
                  }`}>
                    {['converting', 'syncing', 'success'].includes(processingStage) ? (
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    ) : processingStage === 'fetching' ? (
                      <div className="w-2 h-2 bg-indigo-600 rounded-full animate-pulse"></div>
                    ) : (
                      <div className="w-2 h-2 bg-gray-300 rounded-full"></div>
                    )}
                  </div>
                  <span>Fetching content</span>
                </div>
                
                <div className={`flex items-center gap-2 ${
                  ['syncing', 'success'].includes(processingStage) 
                    ? 'text-indigo-600' : 'text-gray-400'
                }`}>
                  <div className={`w-5 h-5 rounded-full flex items-center justify-center ${
                    ['syncing', 'success'].includes(processingStage)
                      ? 'bg-indigo-100' : 'bg-gray-100'
                  }`}>
                    {['success'].includes(processingStage) ? (
                      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    ) : processingStage === 'syncing' ? (
                      <div className="w-2 h-2 bg-indigo-600 rounded-full animate-pulse"></div>
                    ) : (
                      <div className="w-2 h-2 bg-gray-300 rounded-full"></div>
                    )}
                  </div>
                  <span>Syncing</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Manual Paste Modal for WeChat */}
      {manualURL && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-2xl p-8 w-[500px] max-w-[90vw]">
            <div className="text-center mb-4">
              <div className="w-16 h-16 mx-auto bg-orange-100 rounded-full flex items-center justify-center mb-4">
                <svg className="w-8 h-8 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-800 mb-2">Manual Processing Required</h3>
              <p className="text-gray-500 text-sm mb-4">
                This WeChat article requires manual processing. Please copy the article content manually.
              </p>
              <a 
                href={manualURL} 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-indigo-600 hover:text-indigo-700 text-sm underline mb-4 block"
              >
                Open article in browser
              </a>
            </div>
            
            <textarea
              id="manual-content"
              placeholder="Paste the article content here..."
              rows={10}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none mb-4"
            />
            
            <div className="flex gap-3">
              <button
                onClick={() => {
                  setManualURL(null);
                  setInputURL('');
                }}
                className="flex-1 py-2 px-4 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={() => {
                  const content = (document.getElementById('manual-content') as HTMLTextAreaElement)?.value;
                  if (content) {
                    handleManualPaste(content);
                  }
                }}
                className="flex-1 py-2 px-4 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors"
              >
                Send Content
              </button>
            </div>
          </div>
        </div>
      )}

      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-800">OB Sync</h1>
              <p className="text-sm text-gray-500">User: {userID.substring(0, 8)}...</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={downloadUserID}
              className="px-3 py-2 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition-colors"
            >
              Download ID
            </button>
            <button
              onClick={handleLogout}
              className="px-3 py-2 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="flex-1 max-w-4xl mx-auto w-full px-4 py-6 flex flex-col">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-800">Messages</h2>
          <button
            onClick={handleSync}
            disabled={isSyncing}
            className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {isSyncing ? (
                <path className="animate-spin" strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              )}
            </svg>
            {isSyncing ? 'Syncing...' : 'Sync'}
          </button>
        </div>

        <div className="flex-1 bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden flex flex-col">
          <div className="flex-1 overflow-y-auto p-4 space-y-3">
            {messages.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-gray-400">
                <svg className="w-16 h-16 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
                <p>No messages yet</p>
                <p className="text-sm">Send a message or sync to get started</p>
              </div>
            ) : (
              messages.map((msg) => (
                <div
                  key={msg.id}
                  className="p-4 bg-gray-50 rounded-lg border border-gray-100"
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="px-2 py-1 bg-indigo-100 text-indigo-600 text-xs font-medium rounded-full">
                      {getTypeLabel(msg.type)}
                    </span>
                    <span className="text-xs text-gray-400">{formatTime(msg.created_at)}</span>
                  </div>
                  {msg.type === 'attachment' && msg.attachment ? (
                    <div className="flex items-center gap-2">
                      <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      <span className="text-gray-700">{msg.attachment.filename}</span>
                    </div>
                  ) : msg.type === 'url' ? (
                    <a
                      href={msg.original_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-indigo-600 hover:text-indigo-700 hover:underline break-all"
                    >
                      {msg.original_url}
                    </a>
                  ) : (
                    <p className="text-gray-700 whitespace-pre-wrap">{msg.content}</p>
                  )}
                </div>
              ))
            )}
          </div>

          <div className="border-t border-gray-200 p-4">
            {sendStatus === 'success' && !showProcessingOverlay && (
              <div className="mb-3 p-2 bg-green-50 border border-green-200 rounded-lg text-green-600 text-sm text-center">
                Message sent successfully!
              </div>
            )}
            {sendStatus === 'error' && !showProcessingOverlay && (
              <div className="mb-3 p-2 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm text-center">
                Failed to send message
              </div>
            )}
            
            <div className="space-y-3">
              <textarea
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
                placeholder="Enter your message..."
                rows={3}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
              />
              
              <input
                type="text"
                value={inputURL}
                onChange={(e) => setInputURL(e.target.value)}
                placeholder="Enter a URL (optional)"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
              
              <div className="flex items-center gap-3">
                <label className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg cursor-pointer hover:bg-gray-50 transition-colors">
                  <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                  <span className="text-sm text-gray-600">
                    {uploadFile ? uploadFile.name : 'Upload File'}
                  </span>
                  <input
                    type="file"
                    onChange={handleFileChange}
                    className="hidden"
                  />
                </label>
                
                <button
                  onClick={handleSend}
                  disabled={isSending || (!inputText.trim() && !inputURL.trim() && !uploadFile)}
                  className="flex-1 py-3 px-4 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
                >
                  {isSending ? (
                    <>
                      <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Sending...
                    </>
                  ) : (
                    <>
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                      </svg>
                      Send
                    </>
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
