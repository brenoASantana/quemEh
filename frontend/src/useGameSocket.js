import { useCallback, useRef, useState } from 'react';

export const useGameSocket = () => {
  const [gameState, setGameState] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef(null);

  const connect = useCallback((name, room) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?room=${room}&name=${name}`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === 'GAME_STATE') {
          setGameState(msg.payload);
        }
      } catch (err) {
        console.error("Erro ao parsear:", err);
      }
    };

    ws.onerror = (error) => {
      console.error("Erro WebSocket:", error);
    };

    ws.onclose = () => {
      setIsConnected(false);
      setGameState(null);
    };

    socketRef.current = ws;
  }, []);

  const sendMessage = (type, payload) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ type, payload }));
    }
  };

  const disconnect = useCallback(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }
  }, []);

  return { connect, isConnected, gameState, sendMessage, disconnect };
};