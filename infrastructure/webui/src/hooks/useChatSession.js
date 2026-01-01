import { useState, useCallback } from 'react';
import { apiRequest } from '../lib/api';

const INITIAL_MESSAGE = {
    role: 'assistant',
    content: 'Hallo! Ich bin dein AI-Assistent. Frag mich etwas über deine Dokumente, z.B. "Wie hoch waren die Serverkosten?"',
    timestamp: new Date()
};

/**
 * Custom hook for chat session state management (SRP: State Logic Only)
 * @returns {{ messages, isLoading, input, setInput, sendMessage, clearChat }}
 */
export function useChatSession() {
    const [messages, setMessages] = useState([INITIAL_MESSAGE]);
    const [input, setInput] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    // Helper function with retry logic for 502/504 errors
    const fetchWithRetry = async (url, options, retryCount = 0) => {
        try {
            return await apiRequest(url, options);
        } catch (error) {
            const status = error.status || error.response?.status;
            const isGatewayError = status === 502 || status === 504;

            if (isGatewayError && retryCount < 1) {
                console.warn(`Gateway error ${status}, retrying in 1s...`);
                await new Promise(resolve => setTimeout(resolve, 1000));
                return fetchWithRetry(url, options, retryCount + 1);
            }
            throw error;
        }
    };

    const sendMessage = useCallback(async (messageContent) => {
        if (!messageContent.trim() || isLoading) return;

        const userMessage = {
            role: 'user',
            content: messageContent,
            timestamp: new Date()
        };

        setMessages(prev => [...prev, userMessage]);
        setInput('');
        setIsLoading(true);

        try {
            const data = await fetchWithRetry(
                `/api/v1/ask?q=${encodeURIComponent(messageContent)}`,
                { method: 'GET' }
            );

            const aiMessage = {
                role: 'assistant',
                content: data.answer,
                sources: data.cited_sources || [],
                confidence: data.confidence,
                timestamp: new Date()
            };

            setMessages(prev => [...prev, aiMessage]);
        } catch (error) {
            console.error('Chat error:', error);
            setMessages(prev => [...prev, {
                role: 'assistant',
                content: 'Entschuldigung, es gab einen Fehler bei der Verarbeitung deiner Anfrage. Bitte versuche es später erneut.',
                isError: true,
                timestamp: new Date()
            }]);
        } finally {
            setIsLoading(false);
        }
    }, [isLoading]);

    const clearChat = useCallback(() => {
        setMessages([INITIAL_MESSAGE]);
    }, []);

    return {
        messages,
        isLoading,
        input,
        setInput,
        sendMessage,
        clearChat
    };
}
