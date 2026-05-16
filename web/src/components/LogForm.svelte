<script>
  import { createEventDispatcher } from 'svelte';
  import { actionFunction } from './store.js';


  const dispatch = createEventDispatcher();

  export let isRecording = false;
  export let isConnected = false; // Флаг подключения

  // Функция для проверки статуса записи
  const fetchIsRecording = () =>
    fetch('/api/rec/status')
      .then(r => r.ok ? r.json() : Promise.reject(`HTTP error: ${r.status}`))
      .then(({message}) => {
        isConnected = true;
        isRecording = (message == "true");
      })
      .catch(err => {
        console.error('Error fetching recording status:', err);
        isConnected = false;
      });

  fetchIsRecording();

  // Функция запуска записи с отправкой POST-запроса
  async function startRecording() {
    if (!isConnected) {
      console.warn('Cannot start recording: not connected to server');
      return;
    }

    try {
      const response = await fetch('/api/rec/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }

      // Обновляем статус записи после успешного запроса
      isRecording = true;
      dispatch('recordingStarted'); // Опционально: отправляем событие
    } catch (err) {
      console.error('Error starting recording:', err);
      // Можно добавить UI-уведомление об ошибке
    }
  }

  // Функция остановки записи с отправкой POST-запроса
  async function stopRecording() {
    if (!isConnected) {
      console.warn('Cannot stop recording: not connected to server');
      return;
    }

    try {
      const response = await fetch('/api/rec/stop', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }

      // Обновляем статус записи после успешного запроса
      isRecording = false;
      dispatch('recordingStopped'); // Опционально: отправляем событие
      
      $actionFunction()


    } catch (err) {
      console.error('Error stopping recording:', err);
      // Можно добавить UI-уведомление об ошибке
    }
  }
</script>


<div class="log-container">
  <h6>Запись лога</h6>
  <div class="d-flex align-items-center gap-5">
    {#if isRecording}
      <span class="text-danger fw-bold">
        Запись…
      </span>
      <div class="indicator recording"></div>
      <button
        class="btn btn-danger"
        on:click={stopRecording}
        title="Остановить запись лога"
      >
        Остановить запись
      </button>
    {:else}
      <button
        class="btn btn-success"
        on:click={startRecording}
        title="Начните запись лога действий"
        disabled={!isConnected}
      >
        Начать запись
      </button>
    {/if}
  </div>
</div>

<style>
  .log-container {
    background-color: #f8f9fa;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 16px;
  }

  h6 {
    margin-bottom: 12px;
    font-size: 16px;
    font-weight: 600;
  }

  .indicator {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    display: inline-block;
  }

  .recording {
    background-color: #D50000;
    animation: pulse-size 1.5s infinite;
  }

  @keyframes pulse-size {
    0%, 100% { transform: scale(1); opacity: 1; }
    50% { transform: scale(1.2); opacity: 0.7; }
  }

  .btn {
    padding: 8px 16px;
    border-radius: 6px;
    border: none;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.3s ease;
  }

  .btn-success {
    background-color: #4CAF50;
    color: white;
  }

  .btn-danger {
    background-color: #D50000;
    color: white;
  }

  .btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  /* Стили для отключённой кнопки */
  .btn:disabled {
    background-color: #6c757d;
    color: #e9ecef;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .text-danger {
    color: #D50000;
  }

  .fw-bold {
    font-weight: 700;
  }

  .text-muted {
    color: #6c757d;
  }

  @media (max-width: 576px) {
    .d-flex {
      flex-direction: column;
      gap: 10px;
    }
    .btn {
      width: 100%;
    }
  }
</style>
