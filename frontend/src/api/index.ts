export interface Message {
  id: number;
  type: string;
  content: string;
  original_url: string;
  file_path: string;
  created_at: string;
  attachment?: {
    filename: string;
    file_type: string;
  };
}

export interface GenerateUserIDResponse {
  user_id: string;
}

export interface ValidateUserIDResponse {
  valid: boolean;
}

export interface SendMessageResponse {
  message: string;
  id: number;
  needs_manual?: boolean;
  original_url?: string;
}

export async function generateUserID(): Promise<GenerateUserIDResponse> {
  const response = await fetch('/api/user/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!response.ok) {
    throw new Error('Failed to generate user ID');
  }
  return response.json();
}

export async function validateUserID(userID: string): Promise<ValidateUserIDResponse> {
  const response = await fetch('/api/user/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ user_id: userID }),
  });
  if (!response.ok) {
    throw new Error('Failed to validate user ID');
  }
  return response.json();
}

export async function sendMessage(userID: string, type: string, content: string, originalURL?: string): Promise<SendMessageResponse> {
  const response = await fetch('/api/message/send', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userID,
      type,
      content,
      original_url: originalURL || '',
    }),
  });
  if (!response.ok) {
    throw new Error('Failed to send message');
  }
  return response.json();
}

export async function uploadAttachment(userID: string, file: File): Promise<SendMessageResponse> {
  const formData = new FormData();
  formData.append('user_id', userID);
  formData.append('file', file);

  const response = await fetch('/api/message/upload', {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error('Failed to upload attachment');
  }
  return response.json();
}

export async function syncMessages(userID: string, lastSyncTime: string): Promise<Message[]> {
  const response = await fetch('/api/message/sync', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      user_id: userID,
      last_sync_time: lastSyncTime,
    }),
  });
  if (!response.ok) {
    throw new Error('Failed to sync messages');
  }
  return response.json();
}
