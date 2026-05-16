<script>
  import { actionFunction} from './store.js';

  let dumps = [];
  let currentPage = 1;
  const logsPerPage = 15;

  // Отдельная функция для загрузки и обработки логов
  async function loadAndProcessLogs() {
    try {
      const response = await fetch('/api/dumps');

      if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
      }

      // Обновляем данные
      dumps = await response.json();

      // Пересчитываем зависимые переменные
      updatePagination();
    } catch (error) {
      console.error('Ошибка при запросе:', error);
    }
  }

  // Функция для пересчёта пагинации
  function updatePagination() {
    // Вычисляем общее количество страниц
    totalPages = Math.ceil(dumps.length / logsPerPage);

    // Получаем логи для текущей страницы
    currentLogs = [...dumps].reverse().slice(
      (currentPage - 1) * logsPerPage,
      currentPage * logsPerPage
    );
  }

   $: $actionFunction = loadAndProcessLogs;

  // Получаем логи при инициализации компонента
  loadAndProcessLogs();

  // Зависимые переменные (пересчитываются при изменении dumps или currentPage)
  let totalPages = 0;
  let currentLogs = [];

  $: if (dumps.length > 0 && currentPage) {
    updatePagination();
  }

  // Функция перехода на страницу
  const goToPage = (page) => {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
    }
  };
</script>

<div>
  <!-- Пагинация вверху -->
  <nav class="mb-3">
    <ul class="pagination justify-content-center flex-wrap">
      <li class="page-item {currentPage === 1 ? 'disabled' : ''}">
        <button
          class="page-link"
          on:click={() => goToPage(currentPage - 1)}
          disabled={currentPage === 1}
        >
          Предыдущая
        </button>
      </li>

      {#each Array(totalPages) as _, i}
        <li class="page-item {i + 1 === currentPage ? 'active' : ''}">
          <button
            class="page-link"
            on:click={() => goToPage(i + 1)}
          >
            {i + 1}
          </button>
        </li>
      {/each}

      <li class="page-item {currentPage === totalPages ? 'disabled' : ''}">
        <button
          class="page-link"
          on:click={() => goToPage(currentPage + 1)}
          disabled={currentPage === totalPages}
        >
          Следующая
        </button>
      </li>
    </ul>
  </nav>

  <!-- Таблица логов -->
  <div class="table-responsive">
    <table class="table table-striped table-hover table-bordered">
      <thead class="table-light">
        <tr>
          <th scope="col">Дата</th>
          <th scope="col">Вес</th>
        </tr>
      </thead>
      <tbody>
        {#if currentLogs.length > 0}
          {#each currentLogs as log}
            <tr>
              <td>{log.path}</td>
              <td>{log.size} Б</td>
            </tr>
          {/each}
        {:else}
          <tr>
            <td colspan="2" class="text-center text-muted py-4">
              Записей пока нет
            </td>
          </tr>
        {/if}
      </tbody>
    </table>
  </div>

  <!-- Пагинация внизу -->
  <nav>
    <ul class="pagination justify-content-center flex-wrap">
      <li class="page-item {currentPage === 1 ? 'disabled' : ''}">
        <button
          class="page-link"
          on:click={() => goToPage(currentPage - 1)}
          disabled={currentPage === 1}
        >
          Предыдущая
        </button>
      </li>

      {#each Array(totalPages) as _, i}
        <li class="page-item {i + 1 === currentPage ? 'active' : ''}">
          <button
            class="page-link"
            on:click={() => goToPage(i + 1)}
          >
            {i + 1}
          </button>
        </li>
      {/each}

      <li class="page-item {currentPage === totalPages ? 'disabled' : ''}">
        <button
          class="page-link"
          on:click={() => goToPage(currentPage + 1)}
          disabled={currentPage === totalPages}
        >
          Следующая
        </button>
      </li>
    </ul>
  </nav>
</div>

<style>
  /* Стили для пагинации */
  .pagination {
    gap: 0.25rem;
  }

  .page-link {
    padding: 0.375rem 0.75rem;
    margin: 0;
  }

  /* Адаптивность для мобильных устройств */
  @media (max-width: 576px) {
    .table-responsive {
      border-radius: 8px;
      overflow: hidden;
    }

    .page-link {
      padding: 0.25rem 0.5rem;
      font-size: 0.875rem;
    }

    td, th {
      padding: 0.5rem !important;
    }
  }

  /* Улучшение читаемости на планшетах */
  @media (min-width: 577px) and (max-width: 992px) {
    .page-link {
      padding: 0.4rem 0.6rem;
    }
  }
</style>
