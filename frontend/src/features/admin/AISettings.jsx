import { useEffect, useMemo, useState } from 'react';
import adminApi from '../../services/adminApi';

const AISettings = () => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const [form, setForm] = useState({
    provider: 'ollama',
    cloudflare_api_token: '',
    cloudflare_account_id: '',
    cloudflare_ai_model: '',
    cloudflare_ai_base_url: '',
    ollama_base_url: '',
    ollama_model: '',
  });

  const [flags, setFlags] = useState({
    cloudflare_token_configured: false,
  });

  const isCloudflare = form.provider === 'cloudflare';
  const isOllama = form.provider === 'ollama';

  const providerHint = useMemo(() => {
    if (isCloudflare) return 'Cloudflare Workers AI will be used for question correction and review.';
    return 'Local Ollama will be used for question correction and review.';
  }, [isCloudflare]);

  const loadSettings = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await adminApi.getAISettings();
      setForm((prev) => ({
        ...prev,
        provider: data?.provider || 'ollama',
        cloudflare_account_id: data?.cloudflare?.account_id || '',
        cloudflare_ai_model: data?.cloudflare?.model || '@cf/meta/llama-3.1-8b-instruct',
        cloudflare_ai_base_url: data?.cloudflare?.base_url || '',
        ollama_base_url: data?.ollama?.base_url || 'http://localhost:11434',
        ollama_model: data?.ollama?.model || 'qwen2.5:7b',
      }));
      setFlags({
        cloudflare_token_configured: Boolean(data?.cloudflare?.token_configured),
      });
    } catch (err) {
      setError(err?.message || 'Failed to load AI settings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSettings();
  }, []);

  const onChange = (key, value) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setSuccess('');

    try {
      const payload = {
        provider: form.provider,
      };

      if (isCloudflare) {
        payload.cloudflare_account_id = form.cloudflare_account_id;
        payload.cloudflare_ai_model = form.cloudflare_ai_model;
        payload.cloudflare_ai_base_url = form.cloudflare_ai_base_url;
        if (form.cloudflare_api_token.trim()) {
          payload.cloudflare_api_token = form.cloudflare_api_token;
        }
      }

      if (isOllama) {
        payload.ollama_base_url = form.ollama_base_url;
        payload.ollama_model = form.ollama_model;
      }

      await adminApi.updateAISettings(payload);
      setSuccess('AI settings saved successfully.');
      setForm((prev) => ({
        ...prev,
        cloudflare_api_token: '',
      }));
      await loadSettings();
    } catch (err) {
      setError(err?.message || 'Failed to save AI settings');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="p-6 max-w-3xl">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">AI Provider Settings</h1>
        <p className="text-gray-600 mt-1">Choose LLM provider and update credentials directly from admin panel.</p>
      </div>

      {loading && <p className="text-gray-600">Loading settings...</p>}

      {!loading && (
        <form onSubmit={onSubmit} className="space-y-6 bg-white border border-gray-200 rounded-xl p-6 shadow-sm">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Provider</label>
            <select
              value={form.provider}
              onChange={(e) => onChange('provider', e.target.value)}
              className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="ollama">Ollama (Local)</option>
              <option value="cloudflare">Cloudflare Workers AI</option>
            </select>
            <p className="text-xs text-gray-500 mt-2">{providerHint}</p>
          </div>

          {isCloudflare && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Cloudflare Account ID</label>
                <input
                  type="text"
                  value={form.cloudflare_account_id}
                  onChange={(e) => onChange('cloudflare_account_id', e.target.value)}
                  placeholder="your_account_id"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Cloudflare API Token</label>
                <input
                  type="password"
                  value={form.cloudflare_api_token}
                  onChange={(e) => onChange('cloudflare_api_token', e.target.value)}
                  placeholder="your_cf_api_token"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <p className="text-xs text-gray-500 mt-1">
                  {flags.cloudflare_token_configured
                    ? 'A token is already configured. Leave blank to keep existing token.'
                    : 'No token configured yet.'}
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Cloudflare Model</label>
                <input
                  type="text"
                  value={form.cloudflare_ai_model}
                  onChange={(e) => onChange('cloudflare_ai_model', e.target.value)}
                  placeholder="@cf/meta/llama-3.1-8b-instruct"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Cloudflare AI Base URL (optional)</label>
                <input
                  type="text"
                  value={form.cloudflare_ai_base_url}
                  onChange={(e) => onChange('cloudflare_ai_base_url', e.target.value)}
                  placeholder="https://api.cloudflare.com/client/v4/accounts/{account_id}/ai"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </>
          )}

          {isOllama && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Ollama Base URL</label>
                <input
                  type="text"
                  value={form.ollama_base_url}
                  onChange={(e) => onChange('ollama_base_url', e.target.value)}
                  placeholder="http://localhost:11434"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Ollama Model</label>
                <input
                  type="text"
                  value={form.ollama_model}
                  onChange={(e) => onChange('ollama_model', e.target.value)}
                  placeholder="qwen2.5:7b"
                  className="w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </>
          )}

          {error && <p className="text-sm text-red-600">{error}</p>}
          {success && <p className="text-sm text-green-600">{success}</p>}

          <div className="flex justify-end">
            <button
              type="submit"
              disabled={saving}
              className="inline-flex items-center px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60"
            >
              {saving ? 'Saving...' : 'Save AI Settings'}
            </button>
          </div>
        </form>
      )}
    </div>
  );
};

export default AISettings;
