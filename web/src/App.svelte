<script>
  import LogForm from './components/LogForm.svelte';
  import LogsTable from './components/LogsTable.svelte';

  // localStorage.clear('logs')

  let logs = JSON.parse(localStorage.getItem('logs')) || [];
  // let logs = [];
  let isRecording = false;
  let hostname = 'Loading...'

  const startRecording = () => {
    isRecording = true;
  };

  const fetchHostname = () =>
  fetch('/api/hostname')
    .then(r => r.ok ? r.json() : Promise.reject(`HTTP error: ${r.status}`))
    .then(({ Message }) => hostname = Message)
    .catch(() => hostname = 'Error loading');

  fetchHostname();  

  const stopRecording = () => {
    const newLog = {
      id: Date.now(),
      date: new Date().toLocaleString(),
      weight: Math.random() * 100 + 50 // случайный вес для примера
    };
    logs = [newLog, ...logs];
    localStorage.setItem('logs', JSON.stringify(logs));
    isRecording = false;
  };

  const deleteLog = (id) => {
    logs = logs.filter(log => log.id !== id);
    localStorage.setItem('logs', JSON.stringify(logs));
  };
</script>


<header class="header pb-4">
  <div class="container text-center">
    <div class="robot-icon mb-3">
      <i class="fas fa-robot fa-3x text-primary"></i>
    </div>
    <div class="title-and-status">
      <h1 class="display-4 mb-0">{hostname}</h1>
      <div class="online-status mt-2">
        <span class="status-circle"></span>
        Онлайн
      </div>
    </div>
  </div>
</header>

<main class="container mt-4">
  <div class="row">
    <div class="col-12 mb-4">
          <LogForm
            {isRecording}
            on:start={startRecording}
            on:stop={stopRecording}
          />
    </div>
  </div>
  <div class="row">
    <div class="col-12">
      <LogsTable
        {logs}
        on:delete={deleteLog}
      />
    </div>
  </div>
</main>
<style>
  @import 'bootstrap/dist/css/bootstrap.min.css';

  .header {
    background-color: #f8f9fa; /* Светлый фон */
    padding-top: 2rem;
  }

  h1.display-4 {
    font-weight: bold; /* или font-weight: 700; для более жирного шрифта */
  }

  .robot-icon {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  
  .title-and-status {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .online-status {
    display: flex;
    align-items: center;
    font-size: 0.8rem;
    margin-top: 0.5rem; /* Фиксированное расстояние от заголовка */
  }

  .status-circle {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: #28a745;
    margin-right: 8px;
    position: relative; /* Меняем на relative для корректной работы псевдоэлементов */
  }



  .status-circle::before,
  .status-circle::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%); /* Центрируем относительно родительского элемента */
    display: block;
    border-radius: 50%;
    opacity: 0.5;
    border: 1px solid #28a745;
    animation: wave 2s infinite;
  }

  .status-circle::before {
    width: 20px;
    height: 20px;
    animation-delay: 0.5s;
  }

  .status-circle::after {
    width: 30px;
    height: 30px;
  }

  @keyframes wave {
    0% { 
      transform: translate(-50%, -50%) scale(0); 
      opacity: 0.7; 
    }
    50% { 
      transform: translate(-50%, -50%) scale(1); 
      opacity: 0.3; 
    }
    100% { 
      transform: translate(-50%, -50%) scale(0); 
      opacity: 0; 
    }
  }
</style>