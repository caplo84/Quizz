import React from 'react';
import { Pencil, Trash2 } from 'lucide-react';

const TopicList = ({ topics, searchTerm, onEdit, onDeleteRequest }) => {
  const term = searchTerm.trim().toLowerCase();

  const filtered = topics.filter((topic) => {
    if (!term) return true;
    return [topic.name, topic.slug, topic.description]
      .filter(Boolean)
      .some((value) => String(value).toLowerCase().includes(term));
  });

  if (!filtered.length) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-6 text-center text-gray-600 shadow-sm">
        No topics found.
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
      {filtered.map((topic) => (
        <div key={topic.id} className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
          <div className="mb-2 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900">{topic.name}</h3>
            {topic.icon ? (
              <img src={topic.icon} alt={topic.name} className="h-8 w-8 rounded object-cover" />
            ) : null}
          </div>

          <p className="mb-2 text-xs uppercase tracking-wide text-gray-500">{topic.slug}</p>
          <p className="line-clamp-3 min-h-[60px] text-sm text-gray-600">{topic.description || 'No description'}</p>

          <div className="mt-4 flex items-center justify-end gap-2">
            <button
              type="button"
              onClick={() => onEdit(topic)}
              className="inline-flex items-center gap-1 rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-700 hover:bg-gray-100"
            >
              <Pencil className="h-3.5 w-3.5" />
              Edit
            </button>
            <button
              type="button"
              onClick={() => onDeleteRequest(topic.id)}
              className="inline-flex items-center gap-1 rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-700 hover:bg-red-50"
            >
              <Trash2 className="h-3.5 w-3.5" />
              Delete
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default TopicList;
