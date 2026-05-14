<script>
  import { createEventDispatcher } from 'svelte';
  import ConfirmModal from './ConfirmModal.svelte';

  const dispatch = createEventDispatcher();

  export let logs = [];

  let currentPage = 1;
  const logsPerPage = 15;
  let showModal = false;
  let logToDelete = null;

  // Вычисляем общее количество страниц
  $: totalPages = Math.ceil(logs.length / logsPerPage);

  // Получаем логи для текущей страницы
  $: currentLogs = logs.slice(
    (currentPage - 1) * logsPerPage,
    currentPage * logsPerPage
  );

  // Функция перехода на страницу
  const goToPage = (page) => {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
    }
  };

  // Подтверждение удаления лога
  const confirmDelete = (log) => {
    logToDelete = log;
    showModal = true;
  };

  // Обработка удаления лога (подтверждено)
  const handleDelete = () => {
    dispatch('delete', logToDelete.id);
    showModal = false;
    logToDelete = null;
  };

  // Отмена удаления
  const cancelDelete = () => {
    showModal = false;
    logToDelete = null;
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
          <th scope="col" class="text-end">Действия</th>
        </tr>
      </thead>
      <tbody>
        {#if currentLogs.length > 0}
          {#each currentLogs as log (log.id)}
            <tr>
              <td>{log.date}</td>
              <td>{log.weight.toFixed(2)} KB</td>
              <td class="text-end">
                <button
                  class="btn btn-sm btn-outline-danger"
                  on:click={() => confirmDelete(log)}
                  aria-label={`Удалить запись от ${log.date}`}
                >
                  Удалить
                </button>
              </td>
            </tr>
          {/each}
        {:else}
          <tr>
            <td colspan="3" class="text-center text-muted py-4">
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

<!-- Модальное окно подтверждения удаления -->
<ConfirmModal
  show={showModal}
  log={logToDelete}
  on:confirm={handleDelete}
  on:cancel={cancelDelete}
/>

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
