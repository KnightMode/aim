import React from "react";

interface LinkifyProps {
	text: string;
	style?: React.CSSProperties;
	className?: string;
}

/**
 * Component that detects URLs in text and converts them to clickable hyperlinks
 */
export const Linkify: React.FC<LinkifyProps> = ({ text, style, className }) => {
	// Regular expression to detect URLs
	const urlRegex = /(https?:\/\/[^\s]+)/g;

	const renderText = () => {
		const parts = text.split(urlRegex);

		return parts.map((part, index) => {
			// Check if this part is a URL
			if (part.match(urlRegex)) {
				return (
					<a
						key={index}
						href={part}
						target="_blank"
						rel="noopener noreferrer"
						style={{
							color: "#0052CC",
							textDecoration: "underline",
							cursor: "pointer",
							wordBreak: "break-all",
						}}
						onClick={(e) => e.stopPropagation()}
					>
						{part}
					</a>
				);
			}
			// Return plain text
			return <span key={index}>{part}</span>;
		});
	};

	return (
		<span style={style} className={className}>
			{renderText()}
		</span>
	);
};
