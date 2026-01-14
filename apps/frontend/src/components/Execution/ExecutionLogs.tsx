import { useTaskLogs } from "../../hooks/useExecutionLogs";
import { Linkify } from "../Common/Linkify";

interface ExecutionLogsProps {
	taskId: number;
}

const getLogLevelColor = (level: string) => {
	const colors: Record<string, string> = {
		info: "#0065FF",
		success: "#00875A",
		warning: "#FF8B00",
		error: "#DE350B",
	};
	return colors[level] || "#5E6C84";
};

const getLogLevelIcon = (level: string) => {
	const icons: Record<string, string> = {
		info: "◆",
		success: "✓",
		warning: "⚠",
		error: "✕",
	};
	return icons[level] || "•";
};

export const ExecutionLogs = ({ taskId }: ExecutionLogsProps) => {
	const { data: logs = [], isLoading, error } = useTaskLogs(taskId);

	if (isLoading) {
		return (
			<div
				style={{
					background: "#0d1117",
					borderRadius: "6px",
					padding: "24px",
					border: "1px solid #30363d",
					textAlign: "center",
				}}
			>
				<div
					style={{
						display: "inline-flex",
						alignItems: "center",
						gap: "8px",
						color: "#8b949e",
						fontSize: "13px",
						fontFamily: "monospace",
					}}
				>
					<div
						className="loading-spinner"
						style={{ borderTopColor: "#58a6ff" }}
					/>
					Loading execution logs...
				</div>
			</div>
		);
	}

	if (error) {
		return (
			<div
				style={{
					background: "#0d1117",
					borderRadius: "6px",
					padding: "16px",
					border: "1px solid #f85149",
					textAlign: "center",
				}}
			>
				<span
					style={{
						color: "#f85149",
						fontSize: "13px",
						fontFamily: "monospace",
					}}
				>
					✕ Error loading logs
				</span>
			</div>
		);
	}

	if (logs.length === 0) {
		return (
			<div
				style={{
					background: "#0d1117",
					borderRadius: "6px",
					padding: "24px",
					border: "1px solid #30363d",
					textAlign: "center",
				}}
			>
				<div style={{ fontSize: "32px", opacity: 0.3, marginBottom: "8px" }}>
					○
				</div>
				<span
					style={{
						color: "#8b949e",
						fontSize: "13px",
						fontFamily: "monospace",
					}}
				>
					No execution logs yet
				</span>
			</div>
		);
	}

	return (
		<div
			style={{
				background: "#0d1117",
				border: "1px solid #30363d",
				borderRadius: "6px",
				overflow: "hidden",
			}}
		>
			{/* Terminal Header */}
			<div
				style={{
					background: "#161b22",
					borderBottom: "1px solid #30363d",
					padding: "10px 16px",
					display: "flex",
					alignItems: "center",
					justifyContent: "space-between",
				}}
			>
				<div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
					<div
						style={{
							width: "12px",
							height: "12px",
							background: "#f85149",
							borderRadius: "50%",
						}}
					/>
					<div
						style={{
							width: "12px",
							height: "12px",
							background: "#f0883e",
							borderRadius: "50%",
						}}
					/>
					<div
						style={{
							width: "12px",
							height: "12px",
							background: "#3fb950",
							borderRadius: "50%",
						}}
					/>
					<span
						style={{
							marginLeft: "12px",
							fontSize: "11px",
							color: "#8b949e",
							textTransform: "uppercase",
							letterSpacing: "0.5px",
							fontWeight: 600,
							fontFamily: "monospace",
						}}
					>
						Execution Logs
					</span>
				</div>
				<div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
					<div
						style={{
							width: "6px",
							height: "6px",
							background: "#3fb950",
							borderRadius: "50%",
							animation: "pulse 2s ease-in-out infinite",
						}}
					/>
					<span
						style={{
							fontSize: "11px",
							color: "#8b949e",
							fontFamily: "monospace",
						}}
					>
						{logs.length} entries
					</span>
				</div>
			</div>

			{/* Terminal Body */}
			<div
				style={{
					maxHeight: "400px",
					overflowY: "auto",
					padding: "12px",
					fontFamily:
						'ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace',
					fontSize: "12px",
					lineHeight: 1.6,
					position: "relative",
				}}
			>
				{/* Scanline effect */}
				<div
					style={{
						position: "absolute",
						inset: 0,
						pointerEvents: "none",
						opacity: 0.03,
						background:
							"repeating-linear-gradient(0deg, transparent, transparent 2px, #58a6ff 2px, #58a6ff 4px)",
					}}
				/>

				<div style={{ display: "flex", flexDirection: "column", gap: "4px" }}>
					{logs.map((log, index) => (
						<div
							key={log.id}
							style={{
								display: "flex",
								gap: "12px",
								padding: "6px 8px",
								background: "rgba(110, 118, 129, 0.02)",
								borderRadius: "3px",
								transition: "background 0.15s",
								animation: `fadeIn 0.3s ease-out ${index * 0.03}s both`,
							}}
							onMouseEnter={(e) =>
								(e.currentTarget.style.background = "rgba(110, 118, 129, 0.06)")
							}
							onMouseLeave={(e) =>
								(e.currentTarget.style.background = "rgba(110, 118, 129, 0.02)")
							}
						>
							{/* Timestamp */}
							<span
								style={{
									color: "#6e7681",
									flexShrink: 0,
									userSelect: "none",
									minWidth: "70px",
								}}
							>
								[
								{new Date(log.created_at).toLocaleTimeString("en-US", {
									hour12: false,
								})}
								]
							</span>

							{/* Log Level Icon */}
							<span
								style={{
									flexShrink: 0,
									fontWeight: 700,
									color: getLogLevelColor(log.log_level),
								}}
							>
								{getLogLevelIcon(log.log_level)}
							</span>

							{/* Agent Name */}
							{log.agent_name && (
								<span
									style={{
										color: "#a371f7",
										flexShrink: 0,
									}}
								>
									[{log.agent_name}]
								</span>
							)}

							{/* Message */}
							<span
								style={{
									color: "#c9d1d9",
									flex: 1,
									wordBreak: "break-word",
								}}
							>
								<Linkify text={log.message} />
							</span>
						</div>
					))}
				</div>
			</div>

			<style>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: translateX(-4px);
          }
          to {
            opacity: 1;
            transform: translateX(0);
          }
        }
      `}</style>
		</div>
	);
};
