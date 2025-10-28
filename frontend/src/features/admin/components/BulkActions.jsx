import { Download, RefreshCw, Wand2 } from 'lucide-react';

const BulkActions = ({ running, onRunSync, onDownloadImages, onCorrectQuestions }) => {
  const actions = [
    {
      id: 'sync',
      title: 'Run Full GitHub Sync',
      description: 'Fetch latest quizzes/topics from source repository and update data.',
      button: 'Start Sync',
      icon: RefreshCw,
      onClick: onRunSync,
    },
    {
      id: 'images',
      title: 'Download Topic Images',
      description: 'Download all topic icons/images and refresh local assets.',
      button: 'Download Images',
      icon: Download,
      onClick: onDownloadImages,
    },
    {
      id: 'correction',
      title: 'Correct Questions',
      description: 'Run automated question correction and consistency checks.',
      button: 'Run Corrections',
      icon: Wand2,
      onClick: onCorrectQuestions,
    },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
      {actions.map((action) => {
        const Icon = action.icon;
        const isRunning = Boolean(running[action.id]);

        return (
          <div key={action.id} className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm">
            <div className="mb-3 flex items-center gap-2">
              <div className="rounded-lg bg-gray-100 p-2">
                <Icon className={`h-5 w-5 text-gray-700 ${isRunning ? 'animate-spin' : ''}`} />
              </div>
              <h2 className="text-lg font-semibold text-gray-900">{action.title}</h2>
            </div>

            <p className="mb-4 text-sm text-gray-600">{action.description}</p>

            <button
              type="button"
              disabled={isRunning}
              onClick={action.onClick}
              className="w-full rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-gray-400"
            >
              {isRunning ? 'Running...' : action.button}
            </button>
          </div>
        );
      })}
    </div>
  );
};

export default BulkActions;
