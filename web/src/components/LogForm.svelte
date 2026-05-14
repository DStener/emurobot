<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let isRecording = false;
  let recordTime = '00:00:00'; // Таймер записи

  // Имитация отсчёта времени (в реальном проекте можно заменить на точный таймер)
  let timer;

  function startRecording() {
    isRecording = true;
    recordTime = '00:00:00';
    timer = setInterval(() => {
      const seconds = parseInt(recordTime.substring(6)) + 1;
      const minutes = parseInt(recordTime.substring(3, 5)) + Math.floor(seconds / 60);
      const hours = parseInt(recordTime.substring(0, 2)) + Math.floor(minutes / 60);

      recordTime = 
        String(hours % 24).padStart(2, '0') + ':' +
        String(minutes % 60).padStart(2, '0') + ':' +
        String(seconds % 60).padStart(2, '0');
    }, 1000);
  }

  function stopRecording() {
    isRecording = false;
    clearInterval(timer);
  }
</script>

<div class="log-container">
  <h6>Запись лога</h6>
  <div class="d-flex align-items-center gap-5">
    {#if isRecording}
      <span class="text-danger fw-bold">
        Запись… 
        <small class="ms-2 text-muted">{recordTime}</small>
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
