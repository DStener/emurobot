<script>
  import { createEventDispatcher } from 'svelte';
  import { clickOutside } from './clickOutside.js'; // кастомная директива

  const dispatch = createEventDispatcher();

  export let show = false;
  export let log = null;

  const handleConfirm = () => {
    dispatch('confirm');
  };

  const handleCancel = () => {
    dispatch('cancel');
  };
</script>

<div class="modal-backdrop" class:show={show} use:clickOutside={() => handleCancel()}>
  <div class="modal" class:show={show}>
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">Подтверждение удаления</h5>
        <button type="button" class="btn-close" on:click={handleCancel}></button>
      </div>
      <div class="modal-body">
        Вы уверены, что хотите удалить запись от {log?.date} весом {log?.weight?.toFixed(2)} кг?
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" on:click={handleCancel}>
          Отмена
        </button>
        <button type="button" class="btn btn-danger" on:click={handleConfirm}>
          Удалить
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: none;
    z-index: 1040;
  }

  .modal-backdrop.show {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal {
    background: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transform: scale(0.8);
    opacity: 0;
    transition: all 0.3s ease;
    max-width: 90vw;
    margin: 20px;
  }

  .modal.show {
    transform: scale(1);
    opacity: 1;
  }

  .modal-content {
    display: flex;
    flex-direction: column;
  }
</style>
